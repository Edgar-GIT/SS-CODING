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
	choiceExit   = "0"
)

func printMenu() {
	utils.PrintMenuHeader("SS-CODING Project Manager")
	utils.PrintMenuOption(choiceVerify, "Verify errors and vulnerabilities")
	utils.PrintMenuOption(choiceDeps, "Install / Update dependencies")
	utils.PrintMenuOption(choiceDebug, "Debug mode (live reload)")
	utils.PrintMenuOption(choiceRun, "Run website on localhost")
	utils.PrintMenuOption(choiceExit, "Exit")
	utils.PrintDivider()
	fmt.Println()
}
