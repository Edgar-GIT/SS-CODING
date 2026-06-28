package welcomebot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	welcomeMu      sync.Mutex
	welcomeSession *discordgo.Session
)

func intentSetupHelp() string {
	return `Welcome bot needs the "Server Members Intent" (privileged).

In Discord Developer Portal (for THIS welcome bot app, not the music bot):
  1. https://discord.com/developers/applications
  2. Open your Welcome bot application
  3. Bot → Privileged Gateway Intents
  4. Enable "SERVER MEMBERS INTENT"
  5. Save Changes, then restart the bot`
}

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

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMembers
	session.LogLevel = discordgo.LogInformational

	startLogCapture()
	botLog("Welcome bot starting...")

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
		welcomeSession = session
		return nil
	case <-time.After(8 * time.Second):
		_ = session.Close()
		return fmt.Errorf("websocket: close 4014: Disallowed intent(s)\n\n%s", intentSetupHelp())
	}
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
