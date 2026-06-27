package main

import (
	"fmt"

	"ss-coding/utils"
)

const (
	choiceVerify = "1"
	choiceDeps   = "2"
	choiceDebug  = "3"
	choiceRun    = "4"
	choiceStop   = "5"
	choiceExit   = "0"
)

func printMenu() {
	utils.PrintMenuHeader("SS-CODING Project Manager")
	if utils.DevServerRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● Server running at " + utils.DevServerURL))
		fmt.Println()
	}
	utils.PrintMenuOption(choiceVerify, "Verify errors and vulnerabilities")
	utils.PrintMenuOption(choiceDeps, "Install / Update dependencies")
	utils.PrintMenuOption(choiceDebug, "Debug mode (live reload)")
	utils.PrintMenuOption(choiceRun, "Run website on localhost")
	utils.PrintMenuOption(choiceStop, "Stop website")
	utils.PrintMenuOption(choiceExit, "Exit")
	utils.PrintDivider()
	fmt.Println()
}
