package musicbot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
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

	if err := EnsureMusicDependencies(); err != nil {
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

	registerMusicHandlers(session, cfg)

	if err := session.Open(); err != nil {
		return err
	}

	musicSession = session
	return nil
}

func Stop() error {
	musicMu.Lock()
	defer musicMu.Unlock()

	if musicSession == nil {
		return fmt.Errorf("music bot is not running")
	}

	playersMu.Lock()
	for guildID := range players {
		getPlayer(guildID).stopAll(musicSession, "")
	}
	playersMu.Unlock()

	if err := musicSession.Close(); err != nil {
		return err
	}
	musicSession = nil
	return nil
}

func Running() bool {
	musicMu.Lock()
	defer musicMu.Unlock()
	return musicSession != nil
}
