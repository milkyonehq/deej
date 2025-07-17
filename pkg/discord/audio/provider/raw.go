package provider

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/milkyonehq/deej/pkg/discord/audio/queue"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var (
	rawURLRegex = regexp.MustCompile("(?i)^https?://.*\\.(mp3|wav|flac|aac|ogg|m4a|mp4)(\\?.*)?$")
)

type rawMetadata struct {
	Format struct {
		Duration string            `json:"duration"`
		Tags     map[string]string `json:"tags"`
	} `json:"format"`
}

type Raw struct{}

var _ Provider = &Raw{}

func NewRaw() *Raw {
	return &Raw{}
}

func (r *Raw) Name() string {
	return "raw"
}

func (r *Raw) CanHandle(query string) bool {
	return rawURLRegex.MatchString(query)
}

func (r *Raw) Fetch(query string, requestedBy *discordgo.User) (*queue.Track, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", query)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchTrack, err)
	}

	var data rawMetadata
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseTrackMetadata, err)
	}

	seconds, err := strconv.ParseFloat(data.Format.Duration, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseTrackMetadata, err)
	}
	length := time.Duration(seconds * float64(time.Second))

	track := queue.NewTrack(
		stringOrDefault(data.Format.Tags["title"], "Unknown title"),
		stringOrDefault(data.Format.Tags["artist"], "Unknown artist"),
		query,
		query,
		"",
		length,
		requestedBy,
	)

	return track, nil
}
