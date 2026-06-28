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
	playGen           uint64
	prefetchQuery     string
	startingPlayback  bool
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
	gp.queue = append(gp.queue, Track{Title: query, Query: query})
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

func (gp *GuildPlayer) clearPanel() {
	gp.mu.Lock()
	gp.panelChannelID = ""
	gp.panelMessageID = ""
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
	if botHalted() {
		return fmt.Errorf("connect cancelled")
	}

	gp.mu.Lock()
	if gp.vc != nil && gp.voiceChannelID == channelID {
		vc := gp.vc
		gp.mu.Unlock()
		if vc.Status == discordgo.VoiceConnectionStatusReady {
			return nil
		}
	} else {
		gp.mu.Unlock()
	}

	gp.disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	vc, err := session.ChannelVoiceJoin(ctx, gp.guildID, channelID, false, false)
	if err != nil {
		return err
	}

	if err := sleepOrHalt(3 * time.Second); err != nil {
		_ = vc.Disconnect(ctx)
		return err
	}
	if err := waitVoiceReady(vc, 15*time.Second); err != nil {
		_ = vc.Disconnect(ctx)
		return err
	}

	gp.mu.Lock()
	gp.vc = vc
	gp.voiceChannelID = channelID
	gp.wantVoice = true
	gp.mu.Unlock()
	botLogInfo("Connected to voice channel %s", channelID)
	return nil
}

func (gp *GuildPlayer) disconnect() {
	gp.mu.Lock()
	vc := gp.vc
	gp.vc = nil
	gp.voiceChannelID = ""
	gp.wantVoice = false
	gp.mu.Unlock()
	if vc != nil && vc.Status != discordgo.VoiceConnectionStatusDead {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = vc.Disconnect(ctx)
		cancel()
	}
}

func drainStopSignal(ch <-chan struct{}) {
	select {
	case <-ch:
	default:
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

func (gp *GuildPlayer) queueLen() int {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	return len(gp.queue)
}

func (gp *GuildPlayer) beginPlay(session *discordgo.Session, channelID string) {
	gp.mu.Lock()
	select {
	case gp.stopPlayback <- struct{}{}:
	default:
	}
	if gp.vc != nil {
		gp.vc.Speaking(false)
	}
	gp.playing = false
	gp.playGen++
	gen := gp.playGen
	gp.stopPlayback = make(chan struct{}, 1)
	gp.mu.Unlock()
	botLogInfo("beginPlay: gen=%d queue=%d", gen, gp.queueLen())
	go gp.playNext(session, channelID, gen)
}

func (gp *GuildPlayer) playNext(session *discordgo.Session, channelID string, gen uint64) {
	defer func() {
		if r := recover(); r != nil {
			botLogError("playNext panic gen=%d: %v", gen, r)
		}
	}()

	gp.mu.Lock()
	if gp.playGen != gen {
		gp.mu.Unlock()
		botLogInfo("playNext: gen %d superseded, exiting", gen)
		return
	}
	gp.playing = true
	gp.mu.Unlock()

	defer func() {
		gp.mu.Lock()
		if gp.playGen == gen {
			gp.playing = false
		}
		gp.mu.Unlock()
		go refreshPanel(session, gp.guildID)
	}()

	botLogInfo("playNext: started gen=%d", gen)

	for {
		if botHalted() {
			botLogInfo("playNext: halted gen=%d", gen)
			return
		}

		gp.mu.Lock()
		if gp.playGen != gen {
			gp.mu.Unlock()
			return
		}
		gp.mu.Unlock()

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
			botLogInfo("playNext: queue empty gen=%d", gen)
			if stay {
				sendPlain(session, channelID, "Queue empty. Staying in channel (`!stay off` to leave when idle).")
				return
			}
			gp.disconnect()
			sendPlain(session, channelID, "Queue empty. Leaving voice channel.")
			return
		}
		next := gp.queue[0]
		gp.queue = gp.queue[1:]
		track := next
		trackCopy := track
		gp.current = &trackCopy
		gp.positionSec = 0
		gp.playOffsetSec = 0
		gp.skipVoters = nil
		volume := gp.volume
		gp.mu.Unlock()

		botLogInfo("playNext: track %q file=%q url=%q", track.Title, track.FilePath, track.URL)
		sendNewPanel(session, gp, channelID)

		if trackNeedsDownload(track) {
			query := trackSearchQuery(track)
			sendPlain(session, channelID, "Downloading: **"+query+"**…")
			resolved, err := searchYouTube(query)
			if botHalted() {
				return
			}
			if err != nil || resolved == nil {
				botLogError("Download failed for %s: %v", query, err)
				sendPlain(session, channelID, "Could not download: **"+query+"**")
				gp.mu.Lock()
				gp.current = nil
				gp.mu.Unlock()
				continue
			}
			resolved.Query = query
			track = *resolved
			resolvedCopy := track
			gp.mu.Lock()
			gp.current = &resolvedCopy
			gp.mu.Unlock()
			go refreshPanel(session, gp.guildID)
		}

		if track.FilePath != "" {
			if _, err := os.Stat(track.FilePath); err != nil {
				botLogError("audio file missing %s: %v", track.FilePath, err)
				sendPlain(session, channelID, "Audio file missing for: **"+track.Title+"**")
				gp.mu.Lock()
				gp.current = nil
				gp.mu.Unlock()
				continue
			}
		}

		gp.mu.Lock()
		vc := gp.vc
		gp.mu.Unlock()
		if vc == nil {
			botLogWarn("playNext: no voice connection, reconnecting...")
			if err := gp.connect(session, VoiceChannelID); err != nil {
				botLogError("playNext reconnect failed: %v", err)
				sendPlain(session, channelID, "Not connected to voice channel.")
				gp.mu.Lock()
				gp.current = nil
				gp.mu.Unlock()
				return
			}
			gp.mu.Lock()
			vc = gp.vc
			gp.mu.Unlock()
		}

		go gp.prefetchNextTrack(gen)

		if err := gp.playCurrentTrack(session, vc, track, volume); err != nil {
			if botHalted() {
				return
			}
			botLogError("playback error for %s: %v", track.Title, err)
			sendPlain(session, channelID, "Playback error: "+err.Error())
		}
	}
}

func (gp *GuildPlayer) playCurrentTrack(session *discordgo.Session, vc *discordgo.VoiceConnection, track Track, volume float64) error {
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
		go refreshPanel(session, gp.guildID)
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
	options.RawOutput = true
	if startSec > 0 {
		options.StartTime = int(startSec + 0.5)
	}

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
	drainStopSignal(stop)
	framesSent := 0
	const frameDur = 0.02
	for {
		if botHalted() {
			return nil
		}

		select {
		case <-stop:
			botLogInfo("streamTrack: stop signal (%d frames sent)", framesSent)
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
					return fmt.Errorf("no audio frames encoded")
				}
				botLog("Finished track: %s (%d frames)", track.Title, framesSent)
				return nil
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
			go refreshPanel(dgSession, gp.guildID)
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
	killAllDownloads()
	pkillMediaProcesses()
	gp.stopPlaybackNow()
	gp.setStayInChannel(false)
	gp.clearQueue()
	gp.cleanupCurrentFile()
	cleanupAllDownloads()

	gp.mu.Lock()
	gp.current = nil
	gp.playing = false
	gp.playGen++
	gp.stopPlayback = make(chan struct{}, 1)
	gp.panelChannelID = ""
	gp.panelMessageID = ""
	gp.prefetchQuery = ""
	gp.startingPlayback = false
	gp.mu.Unlock()

	go gp.disconnect()

	if channelID != "" {
		go sendPlain(session, channelID, "Stopped.")
	}
}
