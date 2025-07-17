package player

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/corpix/uarand"
	"github.com/milkyonehq/deej/pkg/discord/audio/provider"
	"github.com/milkyonehq/deej/pkg/discord/audio/queue"
	log "github.com/sirupsen/logrus"
	"io"
	"layeh.com/gopus"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	channels      = 2
	frameSize     = (sampleRate * channels) / 100
	maxDataBytes  = (frameSize * 2) * 2
	sampleRate    = 48000
	defaultVolume = 50
)

var (
	ErrNoMatchingProvider       = errors.New("no matching provider found")
	ErrGetUserVoiceState        = errors.New("failed to get user's voice state")
	ErrJoinUserVoiceChannel     = errors.New("failed to join user's voice channel")
	ErrPlayTrack                = errors.New("failed to play track")
	ErrPipeTrack                = errors.New("failed to pipe track")
	ErrStartTrackPlayback       = errors.New("failed to start track playback")
	ErrSendSpeakingNotification = errors.New("failed to send speaking notification")
	ErrCreateOpusEncoder        = errors.New("failed to create opus encoder")
	ErrReadTrack                = errors.New("failed to read track")
	ErrEncodeTrack              = errors.New("failed to encode track")
	ErrSendSuspendSignal        = errors.New("failed to send suspend signal")
	ErrSendResumeSignal         = errors.New("failed to send resume signal")
)

type playerState int

const (
	stateIdle playerState = iota
	statePlaying
	statePaused
)

type Default struct {
	mutex            sync.Mutex
	waitGroup        sync.WaitGroup
	guildID          string
	session          *discordgo.Session
	providerRegistry *provider.Registry
	queue            *queue.Queue
	status           playerState
	skip             chan bool
	pause            chan bool
	volume           int
}

var _ Player = &Default{}

func NewDefault(guildID string, session *discordgo.Session, providerRegistry *provider.Registry) *Default {

	return &Default{
		guildID:          guildID,
		session:          session,
		providerRegistry: providerRegistry,
		queue:            queue.NewQueue(),
		status:           stateIdle,
		skip:             make(chan bool),
		pause:            make(chan bool),
		volume:           defaultVolume,
	}
}

func (d *Default) Search(query string, requestedBy *discordgo.User) (*queue.Track, error) {
	pvr, ok := d.providerRegistry.FindByQuery(query)
	if !ok {
		return nil, ErrNoMatchingProvider
	}

	track, err := pvr.Fetch(query, requestedBy)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (d *Default) playAudio(track *queue.Track, vc *discordgo.VoiceConnection) (err error) {
	cmd := exec.Command("ffmpeg",
		"-user_agent", uarand.GetRandom(),
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-i", track.URL(),
		"-f", "s16le",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", strconv.Itoa(channels),
		"pipe:1",
	)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrPipeTrack, err)
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("%w: %s", ErrStartTrackPlayback, err)
	}

	defer cmd.Process.Kill()

	if err = vc.Speaking(true); err != nil {
		return fmt.Errorf("%w: %s", ErrSendSpeakingNotification, err)
	}

	defer vc.Speaking(false)

	enc, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateOpusEncoder, err)
	}

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		buf := make([]int16, frameSize*channels)

		if err := binary.Read(out, binary.LittleEndian, &buf); err != nil {
			if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			return fmt.Errorf("%w: %s", ErrReadTrack, err)
		}

		for i := range buf {
			buf[i] = int16(int(buf[i]) * d.volume / 100)
		}

		data, err := enc.Encode(buf, frameSize, maxDataBytes)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrEncodeTrack, err)
		}

		select {
		case <-d.skip:
			return nil
		case <-d.pause:
			if err = cmd.Process.Signal(syscall.SIGSTOP); err != nil {
				return fmt.Errorf("%w: %s", ErrSendSuspendSignal, err)
			}
			<-d.pause
			if err = cmd.Process.Signal(syscall.SIGCONT); err != nil {
				return fmt.Errorf("%w: %s", ErrSendResumeSignal, err)
			}
		case vc.OpusSend <- data:
			<-ticker.C
		}
	}
}

func (d *Default) runPlayback() error {
	d.mutex.Lock()
	d.status = statePlaying
	d.mutex.Unlock()

	defer func() {
		d.mutex.Lock()
		d.status = stateIdle
		d.mutex.Unlock()
	}()

	track := d.queue.Peek()

	vs, err := d.session.State.VoiceState(d.guildID, track.RequestedBy().ID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrGetUserVoiceState, err)
	}

	vc, err := d.session.ChannelVoiceJoin(d.guildID, vs.ChannelID, false, true)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrJoinUserVoiceChannel, err)
	}
	defer vc.Disconnect()

	for !d.queue.IsEmpty() {
		track := d.queue.Peek()

		if err := d.playAudio(track, vc); err != nil {
			d.queue.Pop()
			return fmt.Errorf("%w: %s", ErrPlayTrack, err)
		}

		d.queue.Pop()
	}

	return nil
}

func (d *Default) Play(track *queue.Track) {
	d.queue.Add(track)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.status != stateIdle {
		return
	}

	d.waitGroup.Add(1)

	go func() {
		defer d.waitGroup.Done()

		if err := d.runPlayback(); err != nil {
			log.WithError(err).Errorln("Failed to play track.")
		}
	}()
}

func (d *Default) Stop() {
	d.mutex.Lock()
	if d.status == statePaused {
		d.status = statePlaying
		d.pause <- false
	}
	d.mutex.Unlock()

	if !d.queue.IsEmpty() {
		d.skip <- true
		d.queue.Clear()
	}
	d.waitGroup.Wait()
	close(d.pause)
	close(d.skip)
}

func (d *Default) Skip() {
	d.skip <- true
}

func (d *Default) Queue() *queue.Queue {
	return d.queue
}

func (d *Default) Paused() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.status == statePaused
}

func (d *Default) SetPaused(paused bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.status == stateIdle {
		return
	}

	if paused {
		d.status = statePaused
	} else {
		d.status = statePlaying
	}
	d.pause <- paused
}

func (d *Default) Volume() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.volume
}

func (d *Default) SetVolume(volume int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}

	d.volume = volume
}
