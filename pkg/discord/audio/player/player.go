package player

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/queue"
)

type Player interface {
	Search(query string, requestedBy *discordgo.User) (*queue.Track, error)
	Play(track *queue.Track)
	Stop()
	Skip()
	Queue() *queue.Queue
	Paused() bool
	SetPaused(paused bool)
	Volume() int
	SetVolume(volume int)
}
