package main

import (
	"fmt"
	"strconv"
	"strings"

	"ss-coding/manager/transporter"
	"ss-coding/utils"
)

func devServerPort() int {
	portStr := strings.TrimPrefix(utils.DevServerURL, "http://localhost:")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 3000
	}
	return port
}

func runTransporter() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	if transporter.Running() {
		utils.PrintSuccess("Tunnel already active")
		fmt.Println(utils.BoldHiCyan.Apply("  " + transporter.PublicURL()))
		utils.WaitEnter()
		return
	}

	if err := utils.EnsureDependencies(dir); err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	if !utils.DevServerRunning() {
		utils.PrintInfo("Starting local server...")
		if err := utils.StartDevServer(dir, false); err != nil {
			utils.PrintError(err.Error())
			utils.WaitEnter()
			return
		}
	}

	utils.PrintInfo("Creating ngrok tunnel...")
	url, err := transporter.Start(devServerPort())
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintSuccess("Share this link:")
	fmt.Println(utils.BoldHiCyan.Apply("  " + url))
	utils.PrintInfo("Keep the manager running while others test the site")
	utils.WaitEnter()
}

func stopTransporter() {
	if !transporter.Running() {
		return
	}
	_ = transporter.Stop()
}
