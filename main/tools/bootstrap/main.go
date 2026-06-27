package main

import (
	"fmt"
	"os"

	"ss-coding/discord/music_bot/deps"
)

func main() {
	if err := deps.InstallAll(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Bot dependencies ready")
}
