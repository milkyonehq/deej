package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/configuration"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/bot"
	"github.com/milkyonehq/deej/pkg/discord/command"
	"github.com/milkyonehq/deej/pkg/logger"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func initProviders(providerRegistry *provider.Registry) {
	pvrs := []provider.Provider{
		provider.NewRaw(),
		provider.NewYoutube(),
	}
	var pvrNames []string
	for _, pvr := range pvrs {
		providerRegistry.Register(pvr)
		pvrNames = append(pvrNames, pvr.Name())
	}
	log.WithFields(log.Fields{
		"count":     len(pvrNames),
		"providers": strings.Join(pvrNames, ","),
	}).Infoln("Providers successfully registered.")
}

func initCommands(session *discordgo.Session, commandRegistry *command.Registry, playerRegistry *player.Registry, providerRegistry *provider.Registry) error {
	cmds := []command.Command{
		command.NewClear(playerRegistry, providerRegistry),
		command.NewPause(playerRegistry, providerRegistry),
		command.NewPlay(playerRegistry, providerRegistry),
		command.NewQueue(playerRegistry, providerRegistry),
		command.NewShuffle(playerRegistry, providerRegistry),
		command.NewSkip(playerRegistry, providerRegistry),
		command.NewResume(playerRegistry, providerRegistry),
		command.NewVolume(playerRegistry, providerRegistry),
	}
	var cmdNames []string
	for _, cmd := range cmds {
		if err := commandRegistry.Register(session, cmd); err != nil {
			return err
		}
		cmdNames = append(cmdNames, cmd.Name())
	}
	log.WithFields(log.Fields{
		"count":    len(cmdNames),
		"commands": strings.Join(cmdNames, ","),
	}).Infoln("Commands successfully registered.")

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	logger.Init("info")

	config, err := configuration.New()
	if err != nil {
		log.WithError(err).Fatalln("Failed to load configuration.")
	}
	log.Infoln("Configuration successfully loaded.")

	logger.Init(config.LogLevel)

	log.Infoln("Bot is starting...")

	commandRegistry := command.NewRegistry()
	playerRegistry := player.NewRegistry()
	providerRegistry := provider.NewRegistry()

	b, err := bot.New(config.DiscordBotToken, commandRegistry, playerRegistry, providerRegistry)
	if err != nil {
		log.WithError(err).Fatalln("Failed to create bot.")
	}

	if err = b.Start(); err != nil {
		log.WithError(err).Fatalln("Failed to start bot.")
	}
	defer b.Stop()

	initProviders(providerRegistry)

	if err = initCommands(b.Session(), commandRegistry, playerRegistry, providerRegistry); err != nil {
		log.WithError(err).Fatalln("Failed to initialize commands.")
	}

	log.Infoln("Bot is running. Press CTRL+C to exit.")

	<-ctx.Done()

	log.Infoln("Received shutdown signal. DeeJ bot is shutting down...")
}
