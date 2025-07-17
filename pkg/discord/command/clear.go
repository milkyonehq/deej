package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
)

type Clear struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Clear{}

func NewClear(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Clear {
	return &Clear{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (c *Clear) Name() string {
	return "clear"
}

func (c *Clear) Description() string {
	return "Clear the queue."
}

func (c *Clear) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
	}
}

func (c *Clear) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := c.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, c.providerRegistry)
	})

	if guildPlayer.Queue().IsEmpty() {
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

	guildPlayer.Skip()
	guildPlayer.Queue().Clear()
	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("🗑️ Queue cleared").
			Description("The queue has been cleared.").
			Footer("💡 Tips — Use `/skip` to skip the current track.").
			Color(colorSuccess).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
