package deps

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const libdaveVersion = "v1.1.1/cpp"

func InstallAll() error {
	base, err := resolveMusicBotDir()
	if err != nil {
		return err
	}
	if err := installLibDave(base); err != nil {
		return err
	}
	return nil
}

func resolveMusicBotDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "discord", "music_bot")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			candidate := filepath.Join(dir, "discord", "music_bot")
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("music_bot directory not found")
		}
		dir = parent
	}
}

func installLibDave(baseDir string) error {
	libDir := filepath.Join(baseDir, "libdave")
	includeHeader := filepath.Join(libDir, "include", "dave", "dave.h")
	libFile := filepath.Join(libDir, "lib", libdaveLibName())
	if fileExists(includeHeader) && fileExists(libFile) {
		return nil
	}

	osName, arch, err := platformAsset()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://github.com/discord/libdave/releases/download/%s/libdave-%s-%s-boringssl.zip",
		libdaveVersion, osName, arch,
	)

	zipPath := filepath.Join(libDir, "libdave.zip")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		return err
	}
	if err := downloadFile(url, zipPath); err != nil {
		return err
	}
	defer os.Remove(zipPath)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		switch {
		case strings.HasPrefix(file.Name, "include/dave/"):
			dest := filepath.Join(libDir, file.Name)
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return err
			}
			if err := extractZipFile(file, dest); err != nil {
				return err
			}
		case file.Name == "lib/libdave.so", file.Name == "lib/libdave.dylib", file.Name == "bin/libdave.dll":
			dest := filepath.Join(libDir, "lib", filepath.Base(file.Name))
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return err
			}
			if err := extractZipFile(file, dest); err != nil {
				return err
			}
		}
	}

	if !fileExists(includeHeader) || !fileExists(libFile) {
		return fmt.Errorf("libdave install incomplete")
	}
	return nil
}

func platformAsset() (osName, arch string, err error) {
	switch runtime.GOOS {
	case "linux":
		osName = "Linux"
	case "darwin":
		osName = "macOS"
	case "windows":
		osName = "Windows"
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	switch runtime.GOARCH {
	case "amd64":
		arch = "X64"
	case "arm64":
		arch = "ARM64"
	default:
		return "", "", fmt.Errorf("unsupported arch: %s", runtime.GOARCH)
	}
	return osName, arch, nil
}

func libdaveLibName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libdave.dylib"
	case "windows":
		return "libdave.dll"
	default:
		return "libdave.so"
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s (%d)", url, resp.StatusCode)
	}
	tmp := dest + ".part"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	out.Close()
	return os.Rename(tmp, dest)
}

func extractZipFile(file *zip.File, dest string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, src)
	out.Close()
	return err
}
