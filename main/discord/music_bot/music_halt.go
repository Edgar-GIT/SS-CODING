package musicbot

import (
	"fmt"
	"os/exec"
	"sync/atomic"
	"time"
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

func pkillMediaProcesses() {
	_ = exec.Command("pkill", "-9", "yt-dlp").Run()
	_ = exec.Command("pkill", "-9", "ffmpeg").Run()
}
