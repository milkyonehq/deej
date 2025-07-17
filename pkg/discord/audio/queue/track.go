package queue

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type Track struct {
	title        string
	author       string
	length       time.Duration
	url          string
	webpageURL   string
	thumbnailURL string
	requestedBy  *discordgo.User
}

func NewTrack(title, author, url, webpageURL, thumbnailURL string, length time.Duration, requestedBy *discordgo.User) *Track {
	return &Track{
		title:        title,
		author:       author,
		length:       length,
		url:          url,
		webpageURL:   webpageURL,
		thumbnailURL: thumbnailURL,
		requestedBy:  requestedBy,
	}
}

func (track *Track) Title() string {
	return track.title
}

func (track *Track) Author() string {
	return track.author
}

func (track *Track) Len() time.Duration {
	return track.length
}

func (track *Track) URL() string {
	return track.url
}

func (track *Track) WebpageURL() string {
	return track.webpageURL
}

func (track *Track) ThumbnailURL() string {
	return track.thumbnailURL
}

func (track *Track) RequestedBy() *discordgo.User {
	return track.requestedBy
}
