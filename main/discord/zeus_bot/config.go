package zeusbot

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	commandPrefix = "!"
	defaultZeusID = "1204109640976048148"
)

type ZeusConfig struct {
	Token  string
	ZeusID string
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

func LoadZeusConfig() (ZeusConfig, error) {
	_ = loadEnv()
	zeusID := os.Getenv("ZEUS_ID")
	if zeusID == "" {
		zeusID = defaultZeusID
	}
	cfg := ZeusConfig{
		Token:  os.Getenv("ZEUS_TOKEN"),
		ZeusID: zeusID,
	}
	if cfg.Token == "" {
		return cfg, errMissingToken("ZEUS_TOKEN")
	}
	return cfg, nil
}

func imagePath(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		for _, rel := range []string{
			filepath.Join("resources", "img", "bot_pics", name),
			filepath.Join("main", "resources", "img", "bot_pics", name),
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
	return "", errMissingImage(name)
}

func parseUserID(arg string) string {
	arg = strings.TrimSpace(arg)
	if strings.HasPrefix(arg, "<@") && strings.HasSuffix(arg, ">") {
		arg = strings.TrimPrefix(arg, "<@")
		arg = strings.TrimPrefix(arg, "!")
		arg = strings.TrimSuffix(arg, ">")
	}
	return arg
}

func parsePositiveInt(arg string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(arg))
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

type configError string

func (e configError) Error() string { return string(e) }

func errMissingToken(name string) error {
	return configError("missing " + name + " in .env")
}

func errMissingImage(name string) error {
	return configError("image not found: " + name + " (expected in resources/img/bot_pics/)")
}
