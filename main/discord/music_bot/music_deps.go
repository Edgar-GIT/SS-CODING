package musicbot

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

	"ss-coding/discord/music_bot/deps"
)

func EnsureMusicDependencies() error {
	if err := ensureDirs(); err != nil {
		return err
	}
	if err := deps.InstallAll(); err != nil {
		return err
	}
	if err := ensureYTDlp(); err != nil {
		return err
	}
	if _, err := ffmpegPath(); err != nil {
		return err
	}
	return nil
}

func ffmpegPath() (string, error) {
	if override := os.Getenv("FFMPEG_PATH"); override != "" {
		if _, err := os.Stat(override); err == nil {
			return override, nil
		}
	}
	if path := bundledFFmpegPath(); path != "" {
		return path, nil
	}
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}
	local := localFFmpegName()
	if _, err := os.Stat(local); err == nil {
		return local, nil
	}
	if err := downloadFFmpeg(); err != nil {
		return "", err
	}
	if _, err := os.Stat(local); err != nil {
		return "", fmt.Errorf("ffmpeg not available after install")
	}
	return local, nil
}

func bundledFFmpegPath() string {
	oncePaths.Do(initPaths)
	name := "ffmpeg"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	candidates := []string{
		filepath.Join(tempMusicDir, "ffmpeg", "ffmpeg-8.0-full_build", "bin", name),
	}
	matches, _ := filepath.Glob(filepath.Join(tempMusicDir, "ffmpeg", "*", "bin", name))
	candidates = append(candidates, matches...)
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func ytDlpPath() (string, error) {
	local := filepath.Join(binDir, ytDlpBinaryName())
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

func localFFmpegName() string {
	name := "ffmpeg"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(ffmpegDir, name)
}

func ytDlpBinaryName() string {
	if runtime.GOOS == "windows" {
		return "yt-dlp.exe"
	}
	return "yt-dlp"
}

func ensureYTDlp() error {
	_, err := ytDlpPath()
	return err
}

func downloadYtDlp() error {
	url := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	if runtime.GOOS == "windows" {
		url += ".exe"
	}
	dest := filepath.Join(binDir, ytDlpBinaryName())
	return downloadFile(url, dest, 0o755)
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
	zipPath := filepath.Join(ffmpegDir, "ffmpeg.zip")
	if err := downloadFile(url, zipPath, 0o644); err != nil {
		return err
	}
	defer os.Remove(zipPath)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "/bin/ffmpeg.exe") {
			return extractZipFile(file, localFFmpegName())
		}
	}
	return fmt.Errorf("ffmpeg.exe not found in windows archive")
}

func installFFmpegLinux() error {
	url := "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
	archivePath := filepath.Join(ffmpegDir, "ffmpeg.tar.xz")
	if err := downloadFile(url, archivePath, 0o644); err != nil {
		return err
	}
	cmd := exec.Command("tar", "-xJf", archivePath, "-C", ffmpegDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove(archivePath)

	matches, _ := filepath.Glob(filepath.Join(ffmpegDir, "ffmpeg-*-amd64-static", "ffmpeg"))
	if len(matches) == 0 {
		return fmt.Errorf("ffmpeg binary not found after extraction")
	}
	return os.Rename(matches[0], localFFmpegName())
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
