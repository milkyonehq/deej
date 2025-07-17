package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
)

type Pause struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Pause{}

func NewPause(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Pause {
	return &Pause{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (p *Pause) Name() string {
	return "pause"
}

func (p *Pause) Description() string {
	return "Pause the playback."
}

func (p *Pause) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        p.Name(),
		Description: p.Description(),
	}
}

func (p *Pause) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := p.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, p.providerRegistry)
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

	if guildPlayer.Paused() {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("⏸️ Playback paused").
				Description("The playback is already paused.").
				Footer("💡 Tips — Use `/resume` to resume the playback.").
				Color(colorError).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		return
	}

	guildPlayer.SetPaused(true)
	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("⏸️ Playback paused").
			Description("The playback has been paused.").
			Footer("💡 Tips — Use `/resume` to resume the playback.").
			Color(colorSuccess).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
