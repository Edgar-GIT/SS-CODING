package discord

import (
	musicbot "ss-coding/discord/music_bot"
	welcomebot "ss-coding/discord/welcome_bot"
)

func EnableMusicBot() error  { return musicbot.Enable() }
func StopMusicBot() (string, error) { return musicbot.Stop() }
func MusicBotRunning() bool  { return musicbot.Running() }

func EnableWelcomeBot() error  { return welcomebot.Enable() }
func StopWelcomeBot() error    { return welcomebot.Stop() }
func WelcomeBotRunning() bool  { return welcomebot.Running() }

func StopAllBots() {
	_, _ = StopMusicBot()
	_ = StopWelcomeBot()
}
