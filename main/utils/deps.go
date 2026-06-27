package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func NodeModulesPresent(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "node_modules"))
	return err == nil
}

func InstallDependencies(dir string) error {
	return RunInteractive(dir, "bun", "install")
}

func EnsureDependencies(dir string) error {
	if NodeModulesPresent(dir) {
		return nil
	}
	PrintInfo("Dependencies not found, installing...")
	if err := InstallDependencies(dir); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}
	PrintSuccess("Dependencies installed")
	return nil
}
