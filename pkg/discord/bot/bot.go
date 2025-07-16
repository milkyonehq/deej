package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/handler"
)

var (
	ErrCreateSession = fmt.Errorf("failed to create discord session")
	ErrOpenSession   = fmt.Errorf("failed to open discord session")
	ErrCloseSession  = fmt.Errorf("failed to close discord session")
)

type Bot struct {
	session *discordgo.Session
}

func New(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateSession, err)
	}

	bot := &Bot{
		session: sess,
	}

	bot.session.Identify.Intents = discordgo.IntentsGuildVoiceStates

	interactionCreate := handler.NewInteractionCreate()
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
	if err := b.session.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrCloseSession, err)
	}

	return nil
}

func (b *Bot) Session() *discordgo.Session {
	return b.session
}
