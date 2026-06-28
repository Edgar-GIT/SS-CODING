package discord

import (
	"fmt"
	"strings"
	"time"

	musicbot "ss-coding/discord/music_bot"
	"ss-coding/utils"
)

func printDebugMenu() {
	utils.PrintMenuHeader("Discord Debug Mode")
	printBotStatus()
	utils.PrintMenuOption("1", "Music bot")
	utils.PrintMenuOption("2", "Welcome bot")
	utils.PrintMenuOption("0", "Back")
	utils.PrintDivider()
	fmt.Println()
}

func printMusicDebugMenu() {
	utils.PrintMenuHeader("Music Bot Debug")
	if MusicBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Running"))
		fmt.Println()
	}
	utils.PrintMenuOption("1", "Start bot")
	utils.PrintMenuOption("2", "Stop bot (show logs)")
	utils.PrintMenuOption("0", "Back")
	utils.PrintDivider()
	fmt.Println()
}

func printWelcomeDebugMenu() {
	utils.PrintMenuHeader("Welcome Bot Debug")
	if WelcomeBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Running"))
		fmt.Println()
	}
	utils.PrintMenuOption("1", "Start bot")
	utils.PrintMenuOption("2", "Stop bot (show logs)")
	utils.PrintMenuOption("0", "Back")
	utils.PrintDivider()
	fmt.Println()
}

func startMusicBotDebug() {
	if MusicBotRunning() {
		utils.PrintInfo("Music bot already running")
		return
	}
	utils.PrintInfo("Starting music bot...")
	if err := EnableMusicBot(); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}
	utils.PrintSuccess("Music bot online")
}

func stopMusicBotDebug() {
	utils.PrintInfo("Force stopping music bot...")

	type result struct {
		logs string
		err  error
	}
	done := make(chan result, 1)
	go func() {
		logs, err := StopMusicBot()
		done <- result{logs, err}
	}()

	var logs string
	var err error
	select {
	case r := <-done:
		logs, err = r.logs, r.err
	case <-time.After(2 * time.Second):
		logs = musicbot.SessionLogs()
		if strings.TrimSpace(logs) == "" {
			logs = "  (force stop — bot was not responding, killed downloads)"
		}
		err = fmt.Errorf("stop timed out; downloads were killed")
	}

	utils.ClearTerminal()
	utils.PrintMenuHeader("Music Bot — Session Logs")
	if err != nil {
		utils.PrintError(err.Error())
		fmt.Println()
	}
	if strings.TrimSpace(logs) == "" {
		fmt.Println(utils.Muted("  (no logs captured)"))
	} else {
		fmt.Println(logs)
	}
	utils.PrintDivider()
	utils.WaitEnter()
}

func startWelcomeBotDebug() {
	if WelcomeBotRunning() {
		utils.PrintInfo("Welcome bot already running")
		return
	}
	utils.PrintInfo("Starting welcome bot...")
	if err := EnableWelcomeBot(); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}
	utils.PrintSuccess("Welcome bot online")
}

func stopWelcomeBotDebug() {
	if !WelcomeBotRunning() {
		utils.PrintInfo("Welcome bot is not running")
		utils.WaitEnter()
		return
	}
	utils.PrintInfo("Stopping welcome bot...")
	if err := StopWelcomeBot(); err != nil {
		utils.PrintError(err.Error())
	}
	utils.WaitEnter()
}

func runMusicDebugMenu() {
	for {
		utils.ClearTerminal()
		printMusicDebugMenu()
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			startMusicBotDebug()
		case "2":
			stopMusicBotDebug()
		case "0":
			return
		default:
			utils.PrintError("Invalid option")
		}
	}
}

func runWelcomeDebugMenu() {
	for {
		utils.ClearTerminal()
		printWelcomeDebugMenu()
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			startWelcomeBotDebug()
		case "2":
			stopWelcomeBotDebug()
		case "0":
			return
		default:
			utils.PrintError("Invalid option")
		}
	}
}

func runDebugMenu() {
	for {
		utils.ClearTerminal()
		printDebugMenu()
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			runMusicDebugMenu()
		case "2":
			runWelcomeDebugMenu()
		case "0":
			return
		default:
			utils.PrintError("Invalid option")
		}
	}
}
