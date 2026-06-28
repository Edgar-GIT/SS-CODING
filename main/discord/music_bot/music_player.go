package musicbot

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"ss-coding/discord/deps"
)

type GuildPlayer struct {
	mu             sync.Mutex
	queue          []Track
	history        []Track
	current        *Track
	looping        bool
	volume         float64
	page           int
	vc             *discordgo.VoiceConnection
	voiceChannelID string
	panelChannelID string
	panelMessageID string
	guildID        string
	playing        bool
	paused         bool
	positionSec    float64
	playOffsetSec  float64
	seekInterrupt  bool
	stopPlayback   chan struct{}
	stayInChannel  bool
	wantVoice      bool
	reconnecting   bool
	textChannelID  string
	skipVoters     map[string]struct{}
}

var (
	playersMu sync.Mutex
	players   = map[string]*GuildPlayer{}
)

func getPlayer(guildID string) *GuildPlayer {
	playersMu.Lock()
	defer playersMu.Unlock()
	player, ok := players[guildID]
	if !ok {
		player = &GuildPlayer{
			guildID:      guildID,
			volume:       1.0,
			stopPlayback: make(chan struct{}, 1),
		}
		players[guildID] = player
	}
	return player
}

func (gp *GuildPlayer) enqueue(track Track) {
	gp.mu.Lock()
	gp.queue = append(gp.queue, track)
	gp.mu.Unlock()
}

func (gp *GuildPlayer) enqueueQuery(query string) {
	gp.mu.Lock()
	gp.queue = append(gp.queue, Track{Title: query})
	gp.mu.Unlock()
}

func (gp *GuildPlayer) clearQueue() {
	gp.mu.Lock()
	gp.queue = nil
	gp.mu.Unlock()
}

func (gp *GuildPlayer) setLoop(enabled bool) {
	gp.mu.Lock()
	gp.looping = enabled
	gp.mu.Unlock()
}

func (gp *GuildPlayer) togglePause() bool {
	gp.mu.Lock()
	gp.paused = !gp.paused
	value := gp.paused
	gp.mu.Unlock()
	return value
}

func (gp *GuildPlayer) toggleLoop() bool {
	gp.mu.Lock()
	gp.looping = !gp.looping
	value := gp.looping
	gp.mu.Unlock()
	return value
}

func (gp *GuildPlayer) adjustVolume(delta float64) float64 {
	gp.mu.Lock()
	gp.volume += delta
	if gp.volume < 0 {
		gp.volume = 0
	}
	if gp.volume > 1 {
		gp.volume = 1
	}
	value := gp.volume
	gp.mu.Unlock()
	return value
}

func (gp *GuildPlayer) setVolume(value float64) {
	gp.mu.Lock()
	gp.volume = value
	gp.mu.Unlock()
}

func (gp *GuildPlayer) snapshot() (current *Track, queue []Track, history []Track, looping bool, volume float64, page int) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	if gp.current != nil {
		copyTrack := *gp.current
		current = &copyTrack
	}
	queue = append([]Track(nil), gp.queue...)
	history = append([]Track(nil), gp.history...)
	return current, queue, history, gp.looping, gp.volume, gp.page
}

func (gp *GuildPlayer) setPanel(channelID, messageID string) {
	gp.mu.Lock()
	gp.panelChannelID = channelID
	gp.panelMessageID = messageID
	gp.mu.Unlock()
}

func (gp *GuildPlayer) panelRef() (string, string) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	return gp.panelChannelID, gp.panelMessageID
}

func (gp *GuildPlayer) progressSnapshot() (position float64, duration int) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	if gp.current != nil {
		return gp.positionSec, gp.current.Duration
	}
	return 0, 0
}

func (gp *GuildPlayer) seekBy(delta float64) bool {
	gp.mu.Lock()
	if gp.current == nil || !gp.playing {
		gp.mu.Unlock()
		return false
	}
	dur := float64(gp.current.Duration)
	newPos := gp.positionSec + delta
	if newPos < 0 {
		newPos = 0
	}
	if dur > 0 && newPos >= dur {
		newPos = dur - 0.5
		if newPos < 0 {
			newPos = 0
		}
	}
	gp.positionSec = newPos
	gp.playOffsetSec = newPos
	gp.seekInterrupt = true
	gp.mu.Unlock()
	gp.stopPlaybackNow()
	return true
}

func (gp *GuildPlayer) setPage(page int) {
	gp.mu.Lock()
	if page < 0 {
		page = 0
	}
	gp.page = page
	gp.mu.Unlock()
}

func (gp *GuildPlayer) pageCount(queueLen int) int {
	if queueLen == 0 {
		return 1
	}
	return (queueLen-1)/10 + 1
}

func (gp *GuildPlayer) connect(session *discordgo.Session, channelID string) error {
	gp.mu.Lock()
	if gp.vc != nil && gp.vc.Status == discordgo.VoiceConnectionStatusReady {
		if gp.voiceChannelID == channelID {
			gp.mu.Unlock()
			return nil
		}
		vc := gp.vc
		gp.mu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = vc.Disconnect(ctx)
		cancel()
		gp.mu.Lock()
		gp.vc = nil
		gp.voiceChannelID = ""
	}
	gp.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	vc, err := session.ChannelVoiceJoin(ctx, gp.guildID, channelID, false, false)
	if err != nil {
		return err
	}
	// DAVE sender key is derived after the voice websocket handshake completes.
	time.Sleep(3 * time.Second)
	gp.mu.Lock()
	gp.vc = vc
	gp.voiceChannelID = channelID
	gp.wantVoice = true
	gp.mu.Unlock()
	return nil
}

func (gp *GuildPlayer) disconnect() {
	gp.mu.Lock()
	vc := gp.vc
	gp.vc = nil
	gp.voiceChannelID = ""
	gp.wantVoice = false
	gp.playing = false
	gp.mu.Unlock()
	if vc != nil && vc.Status != discordgo.VoiceConnectionStatusDead {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = vc.Disconnect(ctx)
		cancel()
	}
}

func (gp *GuildPlayer) stopPlaybackNow() {
	select {
	case gp.stopPlayback <- struct{}{}:
	default:
	}
	gp.mu.Lock()
	if gp.vc != nil {
		gp.vc.Speaking(false)
	}
	gp.mu.Unlock()
}

func (gp *GuildPlayer) isPlaying() bool {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	return gp.playing
}

func (gp *GuildPlayer) startIfIdle(session *discordgo.Session, channelID string) {
	if gp.isPlaying() {
		return
	}
	go gp.playNext(session, channelID)
}

func (gp *GuildPlayer) playNext(session *discordgo.Session, channelID string) {
	gp.mu.Lock()
	if gp.playing {
		gp.mu.Unlock()
		return
	}
	gp.playing = true
	gp.mu.Unlock()

	defer func() {
		gp.mu.Lock()
		gp.playing = false
		gp.mu.Unlock()
		refreshPanel(session, gp.guildID)
	}()

	for {
		gp.cleanupCurrentFile()

		gp.mu.Lock()
		if gp.current != nil {
			if gp.looping {
				gp.queue = append([]Track{*gp.current}, gp.queue...)
			} else {
				gp.history = append(gp.history, *gp.current)
			}
		}
		if len(gp.queue) == 0 {
			gp.current = nil
			stay := gp.stayInChannel
			gp.mu.Unlock()
			if stay {
				botLogInfo("Queue empty — staying in voice channel")
				sendPlain(session, channelID, "✅ Queue empty. Staying in channel (`!stay off` to leave when idle).")
				return
			}
			gp.disconnect()
			sendPlain(session, channelID, "✅ Playlist empty. Leaving channel.")
			return
		}
		next := gp.queue[0]
		gp.queue = gp.queue[1:]
		gp.current = &next
		gp.positionSec = 0
		gp.playOffsetSec = 0
		gp.resetSkipVotes()
		vc := gp.vc
		volume := gp.volume
		track := next
		gp.mu.Unlock()

		if track.URL == "" && track.FilePath == "" && track.Title != "" {
			resolved, err := searchYouTube(track.Title)
			if err != nil || resolved == nil {
				sendPlain(session, channelID, "❌ Could not find: "+track.Title)
				gp.mu.Lock()
				gp.current = nil
				gp.mu.Unlock()
				continue
			}
			track = *resolved
			gp.mu.Lock()
			gp.current = &track
			gp.mu.Unlock()
		}

		if vc == nil || vc.Status != discordgo.VoiceConnectionStatusReady {
			botLog("voice not ready (status=%v)", vc.Status)
			sendPlain(session, channelID, "❌ Voice connection not ready.")
			return
		}

		if err := gp.playCurrentTrack(session, channelID, vc, track, volume); err != nil {
			botLog("playback error for %s: %v", track.Title, err)
			sendPlain(session, channelID, "❌ Error playing music: "+err.Error())
		}
	}
}

func (gp *GuildPlayer) playCurrentTrack(session *discordgo.Session, channelID string, vc *discordgo.VoiceConnection, track Track, volume float64) error {
	for {
		gp.mu.Lock()
		offset := gp.playOffsetSec
		gp.mu.Unlock()

		err := gp.streamTrack(session, vc, track, volume, offset)

		gp.mu.Lock()
		interrupted := gp.seekInterrupt
		gp.seekInterrupt = false
		gp.mu.Unlock()

		if !interrupted {
			return err
		}
		refreshPanel(session, gp.guildID)
	}
}

func (gp *GuildPlayer) cleanupCurrentFile() {
	gp.mu.Lock()
	current := gp.current
	gp.mu.Unlock()
	if current == nil || current.FilePath == "" {
		return
	}
	_ = os.Remove(current.FilePath)
}

func (gp *GuildPlayer) streamTrack(dgSession *discordgo.Session, vc *discordgo.VoiceConnection, track Track, volume float64, startSec float64) error {
	input := track.FilePath
	if input == "" {
		input = track.URL
	}
	if input == "" {
		return io.ErrUnexpectedEOF
	}

	ffmpeg, err := deps.FFmpegPath()
	if err != nil {
		return err
	}
	_ = os.Setenv("FFMPEG_PATH", ffmpeg)
	_ = os.Setenv("PATH", filepath.Dir(ffmpeg)+string(os.PathListSeparator)+os.Getenv("PATH"))

	botLog("Playing: %s (input: %s)", track.Title, input)

	options := *dca.StdEncodeOptions
	options.Volume = int(volume * 256)
	options.RawOutput = false
	options.StartTime = int(startSec)

	gp.mu.Lock()
	gp.positionSec = startSec
	gp.mu.Unlock()

	encode, err := dca.EncodeFile(input, &options)
	if err != nil {
		botLog("dca encode error: %v", err)
		return err
	}
	defer encode.Cleanup()

	vc.Speaking(true)
	defer vc.Speaking(false)

	stop := gp.stopPlayback
	framesSent := 0
	const frameDur = 0.02 // 20ms opus frames
	for {
		select {
		case <-stop:
			return nil
		default:
		}

		gp.mu.Lock()
		paused := gp.paused
		gp.mu.Unlock()
		for paused {
			time.Sleep(200 * time.Millisecond)
			gp.mu.Lock()
			paused = gp.paused
			gp.mu.Unlock()
			select {
			case <-stop:
				return nil
			default:
			}
		}

		frame, err := encode.OpusFrame()
		if err != nil {
			if err == io.EOF {
				if framesSent == 0 {
					if msg := encode.FFMPEGMessages(); msg != "" {
						botLog("ffmpeg output:\n%s", msg)
					}
					if encErr := encode.Error(); encErr != nil {
						botLog("encode error: %v", encErr)
					}
					return fmt.Errorf("no audio frames encoded")
				}
				botLog("Finished track: %s (%d frames)", track.Title, framesSent)
				return nil
			}
			if msg := encode.FFMPEGMessages(); msg != "" {
				botLog("ffmpeg output:\n%s", msg)
			}
			botLog("opus frame error: %v", err)
			return err
		}
		if frame == nil {
			return fmt.Errorf("empty opus frame")
		}
		framesSent++
		gp.mu.Lock()
		gp.positionSec += frameDur
		gp.mu.Unlock()

		if framesSent%50 == 0 {
			refreshPanel(dgSession, gp.guildID)
		}

		select {
		case vc.OpusSend <- frame:
		case <-time.After(5 * time.Second):
			return io.ErrClosedPipe
		case <-stop:
			return nil
		}
	}
}

func (gp *GuildPlayer) skip(session *discordgo.Session, channelID string) {
	gp.stopPlaybackNow()
}

func (gp *GuildPlayer) playPrevious(session *discordgo.Session, channelID string) bool {
	gp.mu.Lock()
	if len(gp.history) == 0 {
		gp.mu.Unlock()
		return false
	}
	prev := gp.history[len(gp.history)-1]
	gp.history = gp.history[:len(gp.history)-1]
	if gp.current != nil {
		gp.queue = append([]Track{*gp.current}, gp.queue...)
	}
	gp.queue = append([]Track{prev}, gp.queue...)
	gp.mu.Unlock()
	gp.stopPlaybackNow()
	return true
}

func (gp *GuildPlayer) stopAll(session *discordgo.Session, channelID string) {
	gp.setStayInChannel(false)
	gp.clearQueue()
	gp.cleanupCurrentFile()
	gp.mu.Lock()
	gp.current = nil
	gp.mu.Unlock()
	gp.stopPlaybackNow()
	gp.disconnect()
}
