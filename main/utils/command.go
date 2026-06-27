package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

func RunInteractive(dir string, name string, args ...string) error {
	if name == "bun" {
		name = BunExecutable()
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func RunCapture(dir string, name string, args ...string) CommandResult {
	if name == "bun" {
		name = BunExecutable()
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

func PrintStep(label string, result CommandResult) bool {
	ok := result.Err == nil
	status := HiGreen.Apply("PASS")
	detail := result.Stdout + result.Stderr
	if !ok {
		status = HiRed.Apply("FAIL")
		if detail == "" {
			detail = result.Err.Error()
		}
	}
	fmt.Printf("  %s  %s\n", status, BoldWhite.Apply(label))
	if detail != "" {
		fmt.Println(HiBlack.Apply(detail))
	}
	return ok
}

func Muted(s string) string {
	return HiBlack.Apply(s)
}