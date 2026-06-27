package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const DevServerURL = "http://localhost:3000"

var (
	devServerMu  sync.Mutex
	devServerCmd *exec.Cmd
)

func DevServerRunning() bool {
	devServerMu.Lock()
	defer devServerMu.Unlock()
	return devServerCmd != nil && devServerCmd.Process != nil
}

func StartDevServer(dir string, openBrowser bool) error {
	devServerMu.Lock()
	if devServerCmd != nil {
		devServerMu.Unlock()
		return fmt.Errorf("dev server already running at %s", DevServerURL)
	}
	devServerMu.Unlock()

	bun := BunExecutable()
	cmd := exec.Command(bun, "run", "dev")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	devServerMu.Lock()
	devServerCmd = cmd
	devServerMu.Unlock()

	go func() {
		_ = cmd.Wait()
		devServerMu.Lock()
		if devServerCmd == cmd {
			devServerCmd = nil
		}
		devServerMu.Unlock()
	}()

	if openBrowser {
		go func() {
			time.Sleep(2 * time.Second)
			if err := OpenBrowser(DevServerURL); err != nil {
				PrintError(fmt.Sprintf("Could not open browser: %v", err))
				PrintInfo(fmt.Sprintf("Open manually: %s", DevServerURL))
			}
		}()
	}

	return nil
}

func StopDevServer() error {
	devServerMu.Lock()
	cmd := devServerCmd
	devServerMu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("no dev server running")
	}

	if err := killProcessTree(cmd); err != nil {
		return err
	}

	devServerMu.Lock()
	if devServerCmd == cmd {
		devServerCmd = nil
	}
	devServerMu.Unlock()

	return nil
}
