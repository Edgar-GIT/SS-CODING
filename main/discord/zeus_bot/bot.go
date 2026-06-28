package zeusbot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	zeusMu      sync.Mutex
	zeusSession *discordgo.Session
)

func intentSetupHelp() string {
	return `Zeus bot needs privileged gateway intents.

In Discord Developer Portal (for THIS zeus bot app):
  1. https://discord.com/developers/applications
  2. Open your Zeus bot application
  3. Bot → Privileged Gateway Intents
  4. Enable "SERVER MEMBERS INTENT"
  5. Enable "MESSAGE CONTENT INTENT"
  6. Save Changes, then restart the bot`
}

func Enable() error {
	zeusMu.Lock()
	defer zeusMu.Unlock()

	if zeusSession != nil {
		return fmt.Errorf("zeus bot already running")
	}

	cfg, err := LoadZeusConfig()
	if err != nil {
		return err
	}
	cfgMu.Lock()
	zeusCfg = cfg
	cfgMu.Unlock()
	setActivated(false)

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
	botLog("Zeus bot starting...")

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
		zeusSession = session
		return nil
	case <-time.After(8 * time.Second):
		_ = session.Close()
		return fmt.Errorf("websocket: close 4014: Disallowed intent(s)\n\n%s", intentSetupHelp())
	}
}

func Stop() error {
	zeusMu.Lock()
	session := zeusSession
	zeusSession = nil
	zeusMu.Unlock()

	if session == nil {
		return fmt.Errorf("zeus bot is not running")
	}

	setActivated(false)
	botLog("Zeus bot stopping")
	_ = session.Close()
	stopLogCapture()
	return nil
}

func Running() bool {
	zeusMu.Lock()
	defer zeusMu.Unlock()
	return zeusSession != nil
}
