package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func BunExecutable() string {
	if path, err := exec.LookPath("bun"); err == nil {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "bun"
	}
	candidate := filepath.Join(home, ".bun", "bin", "bun")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return "bun"
}
