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
	if MainBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Main bot running"))
	}
	if ZeusBotRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Zeus bot running"))
	}
	if MusicBotRunning() || MainBotRunning() || ZeusBotRunning() {
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

	if !MainBotRunning() {
		utils.PrintInfo("Starting main bot...")
		if err := EnableMainBot(); err != nil {
			utils.PrintError(err.Error())
			return
		}
		utils.PrintSuccess("Main bot online")
	} else {
		utils.PrintInfo("Main bot already running")
	}

	if !ZeusBotRunning() {
		utils.PrintInfo("Starting zeus bot...")
		if err := EnableZeusBot(); err != nil {
			utils.PrintError(err.Error())
			return
		}
		utils.PrintSuccess("Zeus bot online")
	} else {
		utils.PrintInfo("Zeus bot already running")
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
			if nested {
				return
			}
			StopAllBots()
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
