package command

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Description() string
	ApplicationCommand() *discordgo.ApplicationCommand
	Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}
