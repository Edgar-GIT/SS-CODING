package main

import (
	"ss-coding/utils"
)

func runDependencies() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Installing dependencies...")
	if err := utils.RunInteractive(dir, "bun", "install"); err != nil {
		utils.PrintError("Install failed")
		utils.WaitEnter()
		return
	}
	utils.PrintSuccess("Dependencies installed")

	utils.PrintInfo("Updating dependencies...")
	if err := utils.RunInteractive(dir, "bun", "update"); err != nil {
		utils.PrintError("Update failed")
		utils.WaitEnter()
		return
	}
	utils.PrintSuccess("Dependencies updated")
	utils.WaitEnter()
}
