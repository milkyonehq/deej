package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/player"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/util"
	log "github.com/sirupsen/logrus"
)

type Volume struct {
	playerRegistry   *player.Registry
	providerRegistry *provider.Registry
}

var _ Command = &Volume{}

func NewVolume(playerRegistry *player.Registry, providerRegistry *provider.Registry) *Volume {
	return &Volume{
		playerRegistry:   playerRegistry,
		providerRegistry: providerRegistry,
	}
}

func (v *Volume) Name() string {
	return "volume"
}

func (v *Volume) Description() string {
	return "Set the volume of the player."
}

func (v *Volume) ApplicationCommand() *discordgo.ApplicationCommand {
	minValue := 0.0

	return &discordgo.ApplicationCommand{
		Name:        v.Name(),
		Description: v.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "volume",
				Description: "Volume level.",
				Required:    false,
				MinValue:    &minValue,
				MaxValue:    100,
			},
		},
	}
}

func (v *Volume) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		log.WithError(err).WithField("guildID", interaction.GuildID).Errorln("Guild not found.")
		return
	}

	guildPlayer := v.playerRegistry.FindOrCreate(guild.ID, func() player.Player {
		log.WithField("guildID", guild.ID).Infoln("Player successfully registered.")
		return player.NewDefault(guild.ID, session, v.providerRegistry)
	})

	options := interaction.ApplicationCommandData().Options

	if len(options) == 1 {
		guildPlayer.SetVolume(int(options[0].IntValue()))
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("🔊 Volume set").
				Description(fmt.Sprintf("The volume has been set to `%d`%%.", options[0].IntValue())).
				Footer("💡 Tips — Use `/volume` to get the current volume.").
				Color(colorSuccess).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
	} else {
		if err := session.InteractionRespond(interaction.Interaction,
			util.NewEmbedBuilder().
				Title("🔊 Volume").
				Description(fmt.Sprintf("The volume is set to `%d`%%.", guildPlayer.Volume())).
				Footer("💡 Tips — Use `/volume` to set the volume.").
				Color(colorSuccess).
				BuildResponse(true),
		); err != nil {
			log.WithError(err).Errorln("Failed to send message.")
		}
	}
}
