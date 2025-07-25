package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/command"
	"github.com/milkyonehq/deej/pkg/discord/handler"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	ErrCreateSession = fmt.Errorf("failed to create discord session")
	ErrOpenSession   = fmt.Errorf("failed to open discord session")
	ErrCloseSession  = fmt.Errorf("failed to close discord session")
)

type Bot struct {
	session          *discordgo.Session
	commandRegistry  *command.Registry
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

func New(token string, commandRegistry *command.Registry, playerRegistry *player.Registry, providerRegistry *provider.Registry) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateSession, err)
	}

	bot := &Bot{
		session:          sess,
		commandRegistry:  commandRegistry,
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}

	bot.session.Identify.Intents = discordgo.IntentsGuildVoiceStates

	interactionCreate := handler.NewInteractionCreate(bot.commandRegistry)
	ready := handler.NewReady()
	bot.session.AddHandler(interactionCreate.Handle)
	bot.session.AddHandler(ready.Handle)

	return bot, nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("%w: %s", ErrOpenSession, err)
	}

	return nil
}

func (b *Bot) Stop() error {
	cmds := b.commandRegistry.Commands()
	var cmdNames []string
	for id, cmd := range cmds {
		if err := b.commandRegistry.Unregister(b.session, id); err != nil {
			return err
		}
		cmdNames = append(cmdNames, cmd.Name())
	}
	log.WithFields(log.Fields{
		"count":    len(cmdNames),
		"commands": strings.Join(cmdNames, ","),
	}).Infoln("Commands successfully unregistered.")

	var playerGuildIDs []string
	for guildID := range b.playerRegistry.Players() {
		b.playerRegistry.Unregister(guildID)
		playerGuildIDs = append(playerGuildIDs, guildID)
	}
	log.WithFields(log.Fields{
		"count":    len(playerGuildIDs),
		"guildIDs": strings.Join(playerGuildIDs, ","),
	}).Infoln("Players successfully unregistered.")

	pvrs := b.providerRegistry.Providers()
	var pvrNames []string
	for _, pvr := range pvrs {
		b.providerRegistry.Unregister(pvr)
		pvrNames = append(pvrNames, pvr.Name())
	}
	log.WithFields(log.Fields{
		"count":     len(pvrNames),
		"providers": strings.Join(pvrNames, ","),
	}).Infoln("Providers successfully unregistered.")

	if err := b.session.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrCloseSession, err)
	}

	return nil
}

func (b *Bot) Session() *discordgo.Session {
	return b.session
}
