package main

import (
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
		case choiceExit:
			utils.PrintInfo("Goodbye")
			return
		default:
			utils.PrintError("Invalid option")
		}
		utils.ClearTerminal()
		utils.PrintBanner()
	}
}
