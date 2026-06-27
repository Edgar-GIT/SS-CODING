package musicbot

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	commandPrefix    = "!"
	CommandChannelID = "1520542863589376060"
	VoiceChannelID   = "1520542424647204924"
)

type MusicConfig struct {
	MusicToken  string
	GeniusToken string
}

func loadEnv() error {
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
