package utils

import "os/exec"

func PrepareProcessGroup(cmd *exec.Cmd) {
	setProcessGroup(cmd)
}

func KillProcessTree(cmd *exec.Cmd) error {
	return killProcessTree(cmd)
}
