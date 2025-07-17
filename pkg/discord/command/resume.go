package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
)

type Resume struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Resume{}

func NewResume(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Resume {
	return &Resume{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (r *Resume) Name() string {
	return "resume"
}

func (r *Resume) Description() string {
	return "Resume the playback."
}

func (r *Resume) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        r.Name(),
		Description: r.Description(),
	}
}

func (r *Resume) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := r.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, r.providerRegistry)
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

	if !guildPlayer.Paused() {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("▶️ Playback resumed").
				Description("The playback is already resumed.").
				Footer("💡 Tips — Use `/pause` to pause the playback.").
				Color(colorError).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		return
	}

	guildPlayer.SetPaused(false)
	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("▶️ Playback resumed").
			Description("The playback has been resumed.").
			Footer("💡 Tips — Use `/pause` to pause the playback at any time.").
			Color(colorSuccess).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
