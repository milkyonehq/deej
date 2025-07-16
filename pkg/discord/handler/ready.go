package handler

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Ready struct{}

func NewReady() *Ready {
	return &Ready{}
}

func (r *Ready) Handle(session *discordgo.Session, _ *discordgo.Ready) {
	if err := session.UpdateListeningStatus("music"); err != nil {
		log.WithError(err).Errorln("Failed to update status.")
	}
}
