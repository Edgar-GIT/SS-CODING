package discord

import (
	"fmt"

	"ss-coding/utils"
)

const discordBanner = `
▄▄▄▄▄▄   ▄▄▄▄▄  ▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄   ▄▄▄▄▄   ▄▄▄▄▄▄▄   ▄▄▄▄▄▄
███▀▀██▄  ███  █████▀▀▀ ███▀▀▀▀▀ ▄███████▄ ███▀▀███▄ ███▀▀██▄
███  ███  ███   ▀████▄  ███      ███   ███ ███▄▄███▀ ███  ███
███  ███  ███     ▀████ ███      ███▄▄▄███ ███▀▀██▄  ███  ███
██████▀  ▄███▄ ███████▀ ▀███████  ▀█████▀  ███  ▀███ ██████▀
`

func printBotStatus() {
	if MusicBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Music bot running"))
	}
	if WelcomeBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Welcome bot running"))
	}
	if MusicBotRunning() || WelcomeBotRunning() {
		fmt.Println()
	}
}

func printDiscordMenu(nested bool) {
	utils.PrintDivider()
	fmt.Println(utils.BoldHiBlue.Apply(discordBanner))
	utils.PrintDivider()
	printBotStatus()
	utils.PrintMenuOption("1", "Launch all bots")
	utils.PrintMenuOption("2", "Stop all bots")
	utils.PrintMenuOption("3", "Debug mode")
	if nested {
		utils.PrintMenuOption("0", "Back to main menu")
	} else {
		utils.PrintMenuOption("0", "Exit")
	}
	utils.PrintDivider()
	fmt.Println()
}

func printDebugMenu() {
	utils.PrintMenuHeader("Discord Debug Mode")
	printBotStatus()
	utils.PrintMenuOption("1", "Start music bot")
	utils.PrintMenuOption("2", "Start welcome bot")
	utils.PrintMenuOption("3", "Stop music bot")
	utils.PrintMenuOption("4", "Stop welcome bot")
	utils.PrintMenuOption("0", "Back")
	utils.PrintDivider()
	fmt.Println()
}

func launchAllBots() {
	if !MusicBotRunning() {
		utils.PrintInfo("Starting music bot...")
		if err := EnableMusicBot(); err != nil {
			utils.PrintError(err.Error())
			return
		}
		utils.PrintSuccess("Music bot online")
	} else {
		utils.PrintInfo("Music bot already running")
	}

	if !WelcomeBotRunning() {
		utils.PrintInfo("Starting welcome bot...")
		if err := EnableWelcomeBot(); err != nil {
			utils.PrintError(err.Error())
			return
		}
		utils.PrintSuccess("Welcome bot online")
	} else {
		utils.PrintInfo("Welcome bot already running")
	}
}

func runDebugMenu() {
	for {
		utils.ClearTerminal()
		printDebugMenu()
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			if MusicBotRunning() {
				utils.PrintInfo("Music bot already running")
			} else {
				utils.PrintInfo("Starting music bot...")
				if err := EnableMusicBot(); err != nil {
					utils.PrintError(err.Error())
				} else {
					utils.PrintSuccess("Music bot online")
				}
			}
			utils.WaitEnter()
		case "2":
			if WelcomeBotRunning() {
				utils.PrintInfo("Welcome bot already running")
			} else {
				utils.PrintInfo("Starting welcome bot...")
				if err := EnableWelcomeBot(); err != nil {
					utils.PrintError(err.Error())
				} else {
					utils.PrintSuccess("Welcome bot online")
				}
			}
			utils.WaitEnter()
		case "3":
			if err := StopMusicBot(); err != nil {
				utils.PrintError(err.Error())
			} else {
				utils.PrintSuccess("Music bot stopped")
			}
			utils.WaitEnter()
		case "4":
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

func RunMenu(nested bool) {
	for {
		utils.ClearTerminal()
		printDiscordMenu(nested)
		switch utils.ReadChoice("Select an option: ") {
		case "1":
			launchAllBots()
			utils.WaitEnter()
		case "2":
			utils.PrintInfo("Stopping all bots...")
			StopAllBots()
			utils.PrintSuccess("All bots stopped")
			utils.WaitEnter()
		case "3":
			runDebugMenu()
		case "0":
			StopAllBots()
			if nested {
				return
			}
			utils.PrintInfo("Goodbye")
			return
		default:
			utils.PrintError("Invalid option")
		}
	}
}

func main() {
	utils.EnableColors()
	RunMenu(false)
}
