package util

import (
	"github.com/bwmarrin/discordgo"
)

type EmbedBuilder struct {
	embed *discordgo.MessageEmbed
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{
		embed: &discordgo.MessageEmbed{},
	}
}

func (b *EmbedBuilder) Title(title string) *EmbedBuilder {
	b.embed.Title = title

	return b
}

func (b *EmbedBuilder) Description(desc string) *EmbedBuilder {
	b.embed.Description = desc

	return b
}

func (b *EmbedBuilder) Color(color int) *EmbedBuilder {
	b.embed.Color = color

	return b
}

func (b *EmbedBuilder) Footer(text string) *EmbedBuilder {
	b.embed.Footer = &discordgo.MessageEmbedFooter{
		Text: text,
	}

	return b
}

func (b *EmbedBuilder) Thumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}

	return b
}

func (b *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	b.embed.Fields = append(b.embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return b
}

func (b *EmbedBuilder) BuildResponse(ephemeral bool) *discordgo.InteractionResponse {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{b.embed},
			Flags:  flags,
		},
	}
}

func (b *EmbedBuilder) BuildResponseEdit() *discordgo.WebhookEdit {
	return &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{b.embed},
	}
}
