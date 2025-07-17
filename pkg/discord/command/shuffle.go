package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
)

type Shuffle struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Shuffle{}

func NewShuffle(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Shuffle {
	return &Shuffle{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (s *Shuffle) Name() string {
	return "shuffle"
}

func (s *Shuffle) Description() string {
	return "Shuffle the queue."
}

func (s *Shuffle) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        s.Name(),
		Description: s.Description(),
	}
}

func (s *Shuffle) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

	queue.Shuffle()
	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("⏭️ Queue shuffled").
			Description("The queue has been shuffled.").
			Footer("💡 Tips — Use `/queue` to display the queue.").
			Color(colorSuccess).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
