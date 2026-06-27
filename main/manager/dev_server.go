package main

import (
	"fmt"

	"ss-coding/manager/transporter"
	"ss-coding/utils"
)

func runDebugMode() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	if utils.DevServerRunning() {
		utils.PrintError("Server already running — use Stop website to shut it down")
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Starting debug server with live reload")
	utils.PrintInfo("Edit source files — the browser refreshes automatically")
	utils.PrintInfo("Use Stop website from the menu when you are done")
	fmt.Println()

	if err := utils.StartDevServer(dir, false); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintSuccess(fmt.Sprintf("Debug server running at %s", utils.DevServerURL))
	utils.WaitEnter()
}

func runWebsite() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	if utils.DevServerRunning() {
		utils.PrintError("Server already running — use Stop website to shut it down")
		utils.WaitEnter()
		return
	}

	if err := utils.EnsureDependencies(dir); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Starting local server...")

	if err := utils.StartDevServer(dir, true); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintSuccess(fmt.Sprintf("Website running at %s", utils.DevServerURL))
	utils.WaitEnter()
}

func runStopWebsite() {
	if !utils.DevServerRunning() && !transporter.Running() {
		utils.PrintError("Nothing to stop")
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Stopping...")
	stopTransporter()
	if utils.DevServerRunning() {
		if err := utils.StopDevServer(); err != nil {
			utils.PrintError(err.Error())
			utils.WaitEnter()
			return
		}
	}

	utils.PrintSuccess("Stopped")
	utils.WaitEnter()
}
