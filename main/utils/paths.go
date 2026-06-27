package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func WebAppRelPath() string {
	return filepath.Join("src", "web", "typescript", "react", "star-coder-academy")
}

func ModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func WebAppDir() (string, error) {
	root, err := ModuleRoot()
	if err != nil {
		return "", err
	}
	path := filepath.Join(root, WebAppRelPath())
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("web app not found at %s", path)
	}
	return path, nil
}
