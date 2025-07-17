package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type Skip struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Skip{}

func NewSkip(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Skip {
	return &Skip{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (s *Skip) Name() string {
	return "skip"
}

func (s *Skip) Description() string {
	return "Skip the current track."
}

func (s *Skip) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        s.Name(),
		Description: s.Description(),
	}
}

func (s *Skip) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := s.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, s.providerRegistry)
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

	track := queue.Peek()

	guildPlayer.Skip()
	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("⏭️ Track skipped").
			Description("The track has been skipped.").
			Color(colorSuccess).
			Footer("💡 Tips — Use `/clear` to clear the queue.").
			AddField("💿 Title", fmt.Sprintf("[%s](%s)", track.Title(), track.WebpageURL()), false).
			AddField("🎤 Author", track.Author(), true).
			AddField("⏱️ Length", fmt.Sprintf("`%s`", track.Len().Round(time.Second).String()), true).
			AddField("🙍‍ Requested by", track.RequestedBy().Username, true).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
