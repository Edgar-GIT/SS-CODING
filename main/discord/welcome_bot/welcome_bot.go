package welcomebot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	welcomeMu      sync.Mutex
	welcomeSession *discordgo.Session
)

func Enable() error {
	welcomeMu.Lock()
	defer welcomeMu.Unlock()

	if welcomeSession != nil {
		return fmt.Errorf("welcome bot already running")
	}

	cfg, err := LoadWelcomeConfig()
	if err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return err
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers
	session.LogLevel = discordgo.LogInformational

	startLogCapture()
	botLog("Welcome bot starting...")

	registerHandlers(session)

	if err := session.Open(); err != nil {
		return err
	}

	welcomeSession = session
	return nil
}

func Stop() error {
	welcomeMu.Lock()
	session := welcomeSession
	welcomeSession = nil
	welcomeMu.Unlock()

	if session == nil {
		return fmt.Errorf("welcome bot is not running")
	}

	botLog("Welcome bot stopping")
	_ = session.Close()
	stopLogCapture()
	return nil
}

func Running() bool {
	welcomeMu.Lock()
	defer welcomeMu.Unlock()
	return welcomeSession != nil
}
