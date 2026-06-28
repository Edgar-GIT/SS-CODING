package discord

import (
	mainbot "ss-coding/discord/main_bot"
	musicbot "ss-coding/discord/music_bot"
)

func EnableMusicBot() error       { return musicbot.Enable() }
func StopMusicBot() (string, error) { return musicbot.Stop() }
func MusicBotRunning() bool       { return musicbot.Running() }

func EnableMainBot() error  { return mainbot.Enable() }
func StopMainBot() error    { return mainbot.Stop() }
func MainBotRunning() bool  { return mainbot.Running() }

func StopAllBots() {
	_, _ = StopMusicBot()
	_ = StopMainBot()
}
