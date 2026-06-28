package musicbot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"ss-coding/discord/deps"
)

var (
	musicMu      sync.Mutex
	musicSession *discordgo.Session
)

func Enable() error {
	musicMu.Lock()
	defer musicMu.Unlock()

	if musicSession != nil {
		return fmt.Errorf("music bot already running")
	}

	if err := deps.InstallAll(); err != nil {
		return err
	}

	cfg, err := LoadMusicConfig()
	if err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + cfg.MusicToken)
	if err != nil {
		return err
	}

	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsMessageContent

	session.LogLevel = discordgo.LogInformational

	startLogCapture()
	botLog("Music bot starting...")

	registerMusicHandlers(session, cfg)

	if err := session.Open(); err != nil {
		return err
	}

	musicSession = session
	return nil
}

func Stop() (string, error) {
	musicMu.Lock()
	session := musicSession
	musicMu.Unlock()

	if session == nil {
		return "", fmt.Errorf("music bot is not running")
	}

	botLog("Music bot stopping...")

	playersMu.Lock()
	active := make([]*GuildPlayer, 0, len(players))
	for _, gp := range players {
		active = append(active, gp)
	}
	playersMu.Unlock()

	for _, gp := range active {
		gp.stopAll(session, "")
	}

	musicMu.Lock()
	defer musicMu.Unlock()
	if musicSession == nil {
		return stopLogCapture(), nil
	}

	if err := musicSession.Close(); err != nil {
		logs := stopLogCapture()
		musicSession = nil
		return logs, err
	}
	musicSession = nil
	return stopLogCapture(), nil
}

func Running() bool {
	musicMu.Lock()
	defer musicMu.Unlock()
	return musicSession != nil
}
