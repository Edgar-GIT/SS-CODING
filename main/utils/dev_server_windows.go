//go:build windows

package utils

import (
	"fmt"
	"os/exec"
)

func setProcessGroup(cmd *exec.Cmd) {}

func killProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(cmd.Process.Pid)).Run()
}
