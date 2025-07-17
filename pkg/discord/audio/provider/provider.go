package provider

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/queue"
)

type Provider interface {
	Name() string
	CanHandle(query string) bool
	Fetch(query string, requestedBy *discordgo.User) (*queue.Track, error)
}
