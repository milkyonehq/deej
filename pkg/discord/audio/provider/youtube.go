package provider

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/queue"
	"os/exec"
	"regexp"
	"time"
)

var (
	youtubeURLRegex = regexp.MustCompile("^(https?://)?((www\\.)?youtube\\.com|youtu\\.be)/.+$")
	urlRegex        = regexp.MustCompile("^https?://.+$")
)

type youtubeMetadata struct {
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	Uploader     string `json:"uploader"`
	Url          string `json:"url"`
	WebpageURL   string `json:"webpage_url"`
	ThumbnailURL string `json:"thumbnail"`
}

type Youtube struct{}

var _ Provider = &Youtube{}

func NewYoutube() *Youtube {
	return &Youtube{}
}

func (y *Youtube) Name() string {
	return "youtube"
}

func (y *Youtube) CanHandle(query string) bool {
	if youtubeURLRegex.MatchString(query) {
		return true
	}

	if urlRegex.MatchString(query) {
		return false
	}

	return true
}

func (y *Youtube) Fetch(query string, requestedBy *discordgo.User) (*queue.Track, error) {
	matched := youtubeURLRegex.MatchString(query)
	if !matched {
		query = fmt.Sprintf("ytsearch:%s", query)
	}

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-j", query)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchTrack, err)
	}

	var data youtubeMetadata
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseTrackMetadata, err)
	}

	length := time.Duration(data.Duration) * time.Second

	track := queue.NewTrack(
		stringOrDefault(data.Title, "Unknown title"),
		stringOrDefault(data.Uploader, "Unknown artist"),
		data.Url,
		data.WebpageURL,
		data.ThumbnailURL,
		length,
		requestedBy,
	)

	return track, nil
}
