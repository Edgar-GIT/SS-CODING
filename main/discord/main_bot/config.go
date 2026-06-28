package mainbot

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	commandPrefix    = "!"
	WelcomeChannelID = "1520532948900643067"
	RulesChannelID   = "1520572468081987835"
	LogChannelID     = "1520570791169757204"
	MusicChannelID   = "1520542863589376060"
)

type MainConfig struct {
	Token   string
	OwnerID string
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

func LoadMainConfig() (MainConfig, error) {
	_ = loadEnv()
	token := os.Getenv("MAIN_TOKEN")
	if token == "" {
		token = os.Getenv("WELCOME_TOKEN")
	}
	cfg := MainConfig{
		Token:   token,
		OwnerID: os.Getenv("MAIN_OWNER_ID"),
	}
	if cfg.Token == "" {
		return cfg, errMissingToken("MAIN_TOKEN")
	}
	return cfg, nil
}

func welcomeImagePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		for _, rel := range []string{
			filepath.Join("resources", "img", "bot_pics", "welcome.png"),
			filepath.Join("main", "resources", "img", "bot_pics", "welcome.png"),
		} {
			p := filepath.Join(dir, rel)
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errMissingImage()
}

type configError string

func (e configError) Error() string { return string(e) }

func errMissingToken(name string) error {
	return configError("missing " + name + " in .env")
}

func errMissingImage() error {
	return configError("welcome.png not found (expected in resources/img/bot_pics/)")
}
