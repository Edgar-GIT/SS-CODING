package discord

import (
	"fmt"

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
	utils.PrintMenuOption("1", "Start bot (live console)")
	utils.PrintMenuOption("2", "Stop bot")
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
	utils.PrintMenuOption("1", "Start bot (live console)")
	utils.PrintMenuOption("2", "Stop bot")
	utils.PrintMenuOption("0", "Back")
	utils.PrintDivider()
	fmt.Println()
}

func runLiveBotSession(title string, running func() bool, start func() error, stop func() error) {
	utils.EnterAltScreen()
	defer utils.LeaveAltScreen()

	utils.PrintMenuHeader(title)
	fmt.Println(utils.HiCyan.Apply("  Bot output below — press Enter to stop and return to menu."))
	utils.PrintDivider()
	fmt.Println()

	if running() {
		utils.PrintInfo("Bot already running")
	} else {
		utils.PrintInfo("Starting bot...")
		if err := start(); err != nil {
			utils.PrintError(err.Error())
			utils.WaitEnter()
			return
		}
		utils.PrintSuccess("Bot online")
		fmt.Println()
	}

	utils.WaitForStop()
	if err := stop(); err != nil {
		utils.PrintError(err.Error())
	} else {
		utils.PrintSuccess("Bot stopped")
	}
}

func runMusicDebugSession() {
	runLiveBotSession(
		"Music Bot — Live Console",
		MusicBotRunning,
		EnableMusicBot,
		StopMusicBot,
	)
}

func runWelcomeDebugSession() {
	runLiveBotSession(
		"Welcome Bot — Live Console",
		WelcomeBotRunning,
		EnableWelcomeBot,
		StopWelcomeBot,
	)
}

func runMusicDebugMenu() {
	for {
		utils.ClearTerminal()
		printMusicDebugMenu()
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			runMusicDebugSession()
		case "2":
			if err := StopMusicBot(); err != nil {
				utils.PrintError(err.Error())
			} else {
				utils.PrintSuccess("Music bot stopped")
			}
			utils.WaitEnter()
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
			runWelcomeDebugSession()
		case "2":
			if err := StopWelcomeBot(); err != nil {
				utils.PrintError(err.Error())
			} else {
				utils.PrintSuccess("Welcome bot stopped")
			}
			utils.WaitEnter()
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
