package main

import (
	"fmt"
	"time"

	"ss-coding/utils"
)

const defaultDevURL = "http://localhost:3000"

func runDebugMode() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Starting debug server with live reload")
	utils.PrintInfo("Edit source files — the browser refreshes automatically")
	utils.PrintInfo("Press Ctrl+C to stop")
	fmt.Println()

	if err := utils.RunInteractive(dir, "bun", "run", "dev"); err != nil {
		utils.PrintError("Debug server stopped")
	}
	utils.WaitEnter()
}

func runWebsite() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintInfo("Starting local server...")
	go func() {
		time.Sleep(2 * time.Second)
		if err := utils.OpenBrowser(defaultDevURL); err != nil {
			utils.PrintError(fmt.Sprintf("Could not open browser: %v", err))
			utils.PrintInfo(fmt.Sprintf("Open manually: %s", defaultDevURL))
		}
	}()

	if err := utils.RunInteractive(dir, "bun", "run", "dev"); err != nil {
		utils.PrintError("Server stopped")
	}
	utils.WaitEnter()
}
