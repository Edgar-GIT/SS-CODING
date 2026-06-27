package musicbot

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

const (
	commandPrefix    = "!"
	CommandChannelID = "1520542863589376060"
	VoiceChannelID   = "1520542424647204924"
)

var (
	baseDir      string
	musicDir     string
	playlistDir  string
	ffmpegDir    string
	binDir       string
	tempMusicDir string
	oncePaths    sync.Once
)

type MusicConfig struct {
	MusicToken  string
	GeniusToken string
}

func initPaths() {
	baseDir = resolveBaseDir()
	musicDir = filepath.Join(baseDir, "music")
	playlistDir = filepath.Join(baseDir, "playlists")
	ffmpegDir = filepath.Join(baseDir, "ffmpeg")
	binDir = filepath.Join(baseDir, "bin")
	tempMusicDir = filepath.Join(baseDir, "temp_music")
}

func resolveBaseDir() string {
	if dir, err := os.Getwd(); err == nil {
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				candidate := filepath.Join(dir, "discord", "music_bot")
				if _, err := os.Stat(candidate); err == nil {
					return candidate
				}
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "music_bot")
	}
	return "."
}

func ensureDirs() error {
	oncePaths.Do(initPaths)
	for _, dir := range []string{musicDir, playlistDir, ffmpegDir, binDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func loadEnv() error {
	oncePaths.Do(initPaths)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil
}

func moduleRoot() (string, error) {
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
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func LoadMusicConfig() (MusicConfig, error) {
	_ = loadEnv()
	cfg := MusicConfig{
		MusicToken:  os.Getenv("MUSIC_TOKEN"),
		GeniusToken: os.Getenv("GENIUS_TOKEN"),
	}
	if cfg.MusicToken == "" {
		return cfg, errMissingToken("MUSIC_TOKEN")
	}
	return cfg, nil
}

type configError string

func (e configError) Error() string { return string(e) }

func errMissingToken(name string) error {
	return configError("missing " + name + " in .env")
}
