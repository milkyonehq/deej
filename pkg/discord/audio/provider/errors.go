package provider

import (
	"errors"
)

var (
	ErrFetchTrack         = errors.New("failed to fetch track")
	ErrParseTrackMetadata = errors.New("failed to parse track metadata")
)
