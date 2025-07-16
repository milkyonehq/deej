package handler

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionCreate struct{}

func NewInteractionCreate() *InteractionCreate {
	return &InteractionCreate{}
}

func (c *InteractionCreate) Handle(_ *discordgo.Session, _ *discordgo.InteractionCreate) {

}
