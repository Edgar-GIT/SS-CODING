package mainbot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	mainMu      sync.Mutex
	mainSession *discordgo.Session
)

func intentSetupHelp() string {
	return `Main bot needs privileged gateway intents.

In Discord Developer Portal (for THIS main bot app):
  1. https://discord.com/developers/applications
  2. Open your Main bot application
  3. Bot → Privileged Gateway Intents
  4. Enable "SERVER MEMBERS INTENT"
  5. Enable "MESSAGE CONTENT INTENT"
  6. Save Changes, then restart the bot`
}

func Enable() error {
	mainMu.Lock()
	defer mainMu.Unlock()

	if mainSession != nil {
		return fmt.Errorf("main bot already running")
	}

	cfg, err := LoadMainConfig()
	if err != nil {
		return err
	}
	_ = cfg

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return err
	}

	session.Identify.Intents = discordgo.IntentGuilds |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildMessages |
		discordgo.IntentMessageContent
	session.LogLevel = discordgo.LogInformational

	startLogCapture()
	botLog("Main bot starting...")

	readyCh := make(chan error, 1)
	session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		select {
		case readyCh <- nil:
		default:
		}
	})

	registerHandlers(session)

	if err := session.Open(); err != nil {
		return fmt.Errorf("%w\n\n%s", err, intentSetupHelp())
	}

	select {
	case err := <-readyCh:
		if err != nil {
			_ = session.Close()
			return err
		}
		mainSession = session
		return nil
	case <-time.After(8 * time.Second):
		_ = session.Close()
		return fmt.Errorf("websocket: close 4014: Disallowed intent(s)\n\n%s", intentSetupHelp())
	}
}

func Stop() error {
	mainMu.Lock()
	session := mainSession
	mainSession = nil
	mainMu.Unlock()

	if session == nil {
		return fmt.Errorf("main bot is not running")
	}

	botLog("Main bot stopping")
	_ = session.Close()
	stopLogCapture()
	return nil
}

func Running() bool {
	mainMu.Lock()
	defer mainMu.Unlock()
	return mainSession != nil
}
