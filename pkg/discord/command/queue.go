package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"
)

type Queue struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Queue{}

func NewQueue(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Queue {
	return &Queue{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (q *Queue) Name() string {
	return "queue"
}

func (q *Queue) Description() string {
	return "Display the queue."
}

func (q *Queue) ApplicationCommand() *discordgo.ApplicationCommand {
	minValue := 1.0

	return &discordgo.ApplicationCommand{
		Name:        q.Name(),
		Description: q.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "page",
				Description: "Page number.",
				Required:    false,
				MinValue:    &minValue,
			},
		},
	}
}

func (q *Queue) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := q.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, q.providerRegistry)
	})

	queue := guildPlayer.Queue()

	if queue.IsEmpty() {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("🗑️ Queue empty").
				Description("The queue is empty.").
				Footer("💡 Tips — Add a track to the queue and try again.").
				Color(colorError).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		return
	}

	options := interaction.ApplicationCommandData().Options

	queueLength := queue.Len()

	var page int
	totalPages := int(math.Ceil(float64(queueLength) / 10))
	if len(options) == 0 {
		page = 0
	} else {
		page = int(options[0].IntValue()) - 1
	}

	if page >= totalPages {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("❌️ Page not found").
				Description("The page does not exist.").
				Footer("💡 Tips — Choose a different page number and try again.").
				Color(colorError).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		return
	}

	start := page * 10
	end := start + 10
	if end > queueLength {
		end = queueLength
	}

	tracks := queue.Tracks()[start+1 : end]

	var lines []string
	for i, track := range tracks {
		lines = append(lines, fmt.Sprintf("%d. 💿 [%s](%s) - %s - %s",
			start+i,
			track.Title(),
			track.WebpageURL(),
			track.Author(),
			track.Len().Round(time.Second).String()))
	}

	currentTrack := queue.Peek()

	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("🎧 Queue").
			Description(fmt.Sprintf("▶️ Now playing:\n"+
				"💿 [%s](%s) - %s - `%s`\n\n"+
				"⏭️ Up Next:\n%s",
				currentTrack.Title(),
				currentTrack.WebpageURL(),
				currentTrack.Author(),
				currentTrack.Len().Round(time.Second).String(),
				strings.Join(lines, "\n"),
			)).
			Footer(fmt.Sprintf("Page %d of %d", page+1, totalPages)).
			Thumbnail(currentTrack.ThumbnailURL()).
			Color(colorSuccess).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
