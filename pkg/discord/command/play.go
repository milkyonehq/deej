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

type Play struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Play{}

func NewPlay(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Play {
	return &Play{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (p *Play) Name() string {
	return "play"
}

func (p *Play) Description() string {
	return "Play a single track."
}

func (p *Play) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        p.Name(),
		Description: p.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "track",
				Description: "Track name or URL.",
				Required:    true,
			},
		},
	}
}

func (p *Play) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	_, err = session.State.VoiceState(guild.ID, interaction.Member.User.ID)
	if err != nil {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("🔇 Voice channel not found").
				Description("The user's voice channel couldn't be found.").
				Footer("💡 Tips — Join a voice channel and try again.").
				Color(colorError).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		return
	}

	guildPlayer := p.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, p.providerRegistry)
	})

	if err := session.InteractionRespond(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("🔍 Searching Track").
			Description("Looking for your track, please wait...").
			Footer("💡 Tips — You can queue multiple songs one after another.").
			Color(colorInfo).
			BuildResponse(true),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
		return
	}

	options := interaction.ApplicationCommandData().Options
	track, err := guildPlayer.Search(options[0].StringValue(), interaction.Member.User)
	if err != nil {
		if _, err := session.InteractionResponseEdit(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("🔍 Track not found").
				Description("The track couldn't be found.").
				Footer("💡 Tips — Choose a different track and try again.").
				Color(colorError).
				BuildResponseEdit(),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
		log.WithError(err).Errorln("Failed to search track.")
		return
	}

	guildPlayer.Play(track)
	if _, err := session.InteractionResponseEdit(interaction.Interaction,
		util.NewEmbedBuilder().
			Title("▶️ Track added to queue").
			Description("The track has been added to the queue.").
			Color(colorSuccess).
			AddField("💿 Title", fmt.Sprintf("[%s](%s)", track.Title(), track.WebpageURL()), false).
			AddField("🎤 Author", track.Author(), true).
			AddField("⏱️ Length", fmt.Sprintf("`%s`", track.Len().Round(time.Second).String()), true).
			AddField("🙍‍ Requested by", track.RequestedBy().Username, true).
			Thumbnail(track.ThumbnailURL()).
			Footer("💡 Tips — Use `/queue` to display queue.").
			BuildResponseEdit(),
	); err != nil {
		log.WithError(err).Errorln("Failed to send message.")
	}
}
