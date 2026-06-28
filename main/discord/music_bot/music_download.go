package musicbot

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"ss-coding/discord/deps"
)

var activeDownloads struct {
	mu   sync.Mutex
	cmds []*exec.Cmd
}

func registerDownloadCmd(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process == nil {
		// Process not started yet; caller should register after Start or use runDownloadCmd
	}
	activeDownloads.mu.Lock()
	activeDownloads.cmds = append(activeDownloads.cmds, cmd)
	activeDownloads.mu.Unlock()
}

func unregisterDownloadCmd(cmd *exec.Cmd) {
	activeDownloads.mu.Lock()
	defer activeDownloads.mu.Unlock()
	for i, c := range activeDownloads.cmds {
		if c == cmd {
			activeDownloads.cmds = append(activeDownloads.cmds[:i], activeDownloads.cmds[i+1:]...)
			return
		}
	}
}

func killAllDownloads() {
	activeDownloads.mu.Lock()
	cmds := activeDownloads.cmds
	activeDownloads.cmds = nil
	activeDownloads.mu.Unlock()
	for _, cmd := range cmds {
		killProcessTree(cmd)
	}
}

func cleanupAllDownloads() int {
	dir := deps.DownloadsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			botLogWarn("cleanup downloads: %v", err)
		}
		return 0
	}
	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := os.Remove(path); err != nil {
			botLogWarn("cleanup: could not remove %s: %v", path, err)
			continue
		}
		removed++
	}
	return removed
}

func prepareCmd(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessTree(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	}
	_ = cmd.Process.Kill()
}

func runDownloadCmd(cmd *exec.Cmd) error {
	prepareCmd(cmd)
	registerDownloadCmd(cmd)
	defer unregisterDownloadCmd(cmd)
	return cmd.Run()
}

// killActiveDownload is kept for compatibility with music_search.go
func killActiveDownload() {
	killAllDownloads()
}
