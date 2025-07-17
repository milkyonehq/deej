package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/command"
)

type InteractionCreate struct {
	registry *command.Registry
}

func NewInteractionCreate(registry *command.Registry) *InteractionCreate {
	return &InteractionCreate{
		registry: registry,
	}
}

func (c *InteractionCreate) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	cmds := c.registry.Commands()
	for _, cmd := range cmds {
		if interaction.ApplicationCommandData().Name == cmd.Name() {
			cmd.Execute(session, interaction)
			return
		}
	}
}
