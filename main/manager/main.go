package main

import (
	"ss-coding/discord"
	"ss-coding/utils"
)

func main() {
	utils.EnableColors()
	utils.ClearTerminal()
	utils.PrintBanner()

	for {
		printMenu()
		switch utils.ReadChoice("Select an option: ") {
		case choiceVerify:
			runVerify()
		case choiceDeps:
			runDependencies()
		case choiceDebug:
			runDebugMode()
		case choiceRun:
			runWebsite()
		case choiceStop:
			runStopWebsite()
		case choiceTransporter:
			runTransporter()
		case choiceDiscord:
			runDiscordMenu()
		case choiceExit:
			stopTransporter()
			if utils.DevServerRunning() {
				utils.PrintInfo("Stopping server before exit...")
				_ = utils.StopDevServer()
			}
			discord.StopAllBots()
			utils.PrintInfo("Goodbye")
			return
		default:
			utils.PrintError("Invalid option")
		}
		utils.ClearTerminal()
		utils.PrintBanner()
	}
}
