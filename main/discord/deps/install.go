package deps

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// InstallAll creates temp_music layout and installs yt-dlp + ffmpeg if missing.
func InstallAll() error {
	if err := EnsureDirs(); err != nil {
		return err
	}
	if err := ensureYTDlp(); err != nil {
		return err
	}
	if _, err := FFmpegPath(); err != nil {
		return err
	}
	return nil
}

func FFmpegPath() (string, error) {
	if override := os.Getenv("FFMPEG_PATH"); override != "" {
		if _, err := os.Stat(override); err == nil {
			return override, nil
		}
	}
	if path := findFFmpeg(); path != "" {
		return path, nil
	}
	if err := downloadFFmpeg(); err != nil {
		return "", err
	}
	local := localFFmpegBinary()
	if _, err := os.Stat(local); err != nil {
		return "", fmt.Errorf("ffmpeg not available after install")
	}
	return local, nil
}

func YTDlpPath() (string, error) {
	local := localYTDlpBinary()
	if _, err := os.Stat(local); err == nil {
		return local, nil
	}
	if path, err := exec.LookPath("yt-dlp"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("yt-dlp.exe"); err == nil {
		return path, nil
	}
	if err := downloadYtDlp(); err != nil {
		return "", err
	}
	if _, err := os.Stat(local); err != nil {
		return "", fmt.Errorf("yt-dlp not available after install")
	}
	return local, nil
}

func findFFmpeg() string {
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}
	local := localFFmpegBinary()
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return ""
}

func localFFmpegBinary() string {
	name := "ffmpeg"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return toolPath(name)
}

func localYTDlpBinary() string {
	name := "yt-dlp"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return toolPath(name)
}

func ensureYTDlp() error {
	_, err := YTDlpPath()
	return err
}

func downloadYtDlp() error {
	url := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	if runtime.GOOS == "windows" {
		url += ".exe"
	}
	return downloadFile(url, localYTDlpBinary(), 0o755)
}

func downloadFFmpeg() error {
	switch runtime.GOOS {
	case "windows":
		return installFFmpegWindows()
	case "linux":
		return installFFmpegLinux()
	case "darwin":
		return installFFmpegDarwin()
	default:
		return fmt.Errorf("unsupported OS for automatic ffmpeg install: %s", runtime.GOOS)
	}
}

func installFFmpegWindows() error {
	url := "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
	zipPath := toolPath("ffmpeg.zip")
	if err := downloadFile(url, zipPath, 0o644); err != nil {
		return err
	}
	defer os.Remove(zipPath)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dest := localFFmpegBinary()
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "/bin/ffmpeg.exe") {
			return extractZipFile(file, dest)
		}
	}
	return fmt.Errorf("ffmpeg.exe not found in windows archive")
}

func installFFmpegLinux() error {
	url := "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
	archivePath := toolPath("ffmpeg.tar.xz")
	if err := downloadFile(url, archivePath, 0o644); err != nil {
		return err
	}
	defer os.Remove(archivePath)

	extractDir := toolPath("ffmpeg-extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return err
	}
	defer os.RemoveAll(extractDir)

	cmd := exec.Command("tar", "-xJf", archivePath, "-C", extractDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	matches, _ := filepath.Glob(filepath.Join(extractDir, "ffmpeg-*-amd64-static", "ffmpeg"))
	if len(matches) == 0 {
		return fmt.Errorf("ffmpeg binary not found after extraction")
	}
	return os.Rename(matches[0], localFFmpegBinary())
}

func installFFmpegDarwin() error {
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("install ffmpeg manually: brew install ffmpeg")
	}
	cmd := exec.Command("brew", "install", "ffmpeg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadFile(url, dest string, mode os.FileMode) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s (%d)", url, resp.StatusCode)
	}
	tmp := dest + ".part"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
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
