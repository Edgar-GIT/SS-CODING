package musicbot

import (
	"fmt"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botHaltedFlag atomic.Bool

func haltBot() {
	botHaltedFlag.Store(true)
}

func resetHalt() {
	botHaltedFlag.Store(false)
}

func botHalted() bool {
	return botHaltedFlag.Load()
}

func sleepOrHalt(d time.Duration) error {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if botHalted() {
			return fmt.Errorf("cancelled")
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func waitVoiceReady(vc *discordgo.VoiceConnection, timeout time.Duration) error {
	if vc == nil {
		return fmt.Errorf("no voice connection")
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if botHalted() {
			return fmt.Errorf("cancelled")
		}
		switch vc.Status {
		case discordgo.VoiceConnectionStatusReady:
			botLogInfo("Voice connection ready")
			return nil
		case discordgo.VoiceConnectionStatusDead:
			if vc.Err != nil {
				return vc.Err
			}
			return fmt.Errorf("voice connection dead")
		}
		time.Sleep(100 * time.Millisecond)
	}
	botLogWarn("Voice not ready after %s (status=%d), trying playback anyway", timeout, vc.Status)
	return nil
}

func pkillMediaProcesses() {
	_ = exec.Command("pkill", "-9", "yt-dlp").Run()
	_ = exec.Command("pkill", "-9", "ffmpeg").Run()
}
