package main

import (
	"fmt"
	"strconv"
	"strings"

	"ss-coding/manager/transporter"
	"ss-coding/utils"
)

const (
	choiceTransporterStart = "1"
	choiceTransporterLogs  = "2"
	choiceTransporterClear = "3"
	choiceTransporterStop  = "4"
	choiceTransporterBack  = "0"
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
	for {
		printTransporterMenu()
		switch utils.ReadChoice("Select an option: ") {
		case choiceTransporterStart:
			startTransporter()
		case choiceTransporterLogs:
			showTransporterLogs()
		case choiceTransporterClear:
			transporter.ClearRequestLogs()
			utils.PrintSuccess("Ngrok IP logs cleared")
			utils.WaitEnter()
		case choiceTransporterStop:
			if transporter.Running() {
				if err := transporter.Stop(); err != nil {
					utils.PrintError(err.Error())
				} else {
					utils.PrintSuccess("Ngrok tunnel stopped")
				}
			} else {
				utils.PrintError("No ngrok tunnel running")
			}
			utils.WaitEnter()
		case choiceTransporterBack:
			return
		default:
			utils.PrintError("Invalid option")
			utils.WaitEnter()
		}
		utils.ClearTerminal()
		utils.PrintBanner()
	}
}

func printTransporterMenu() {
	utils.PrintMenuHeader("Ngrok Transporter")
	if transporter.Running() {
		fmt.Println(utils.HiPurple.Apply("  ● Tunnel active: " + transporter.PublicURL()))
	}
	if transporter.LoggingProxyRunning() {
		fmt.Println(utils.HiGreen.Apply("  ● IP logging active"))
	}
	if transporter.Running() || transporter.LoggingProxyRunning() {
		fmt.Println()
	}
	utils.PrintMenuOption(choiceTransporterStart, "Start server and ngrok tunnel")
	utils.PrintMenuOption(choiceTransporterLogs, "Ngrok IP logging")
	utils.PrintMenuOption(choiceTransporterClear, "Clear ngrok IP logs")
	utils.PrintMenuOption(choiceTransporterStop, "Stop ngrok tunnel")
	utils.PrintMenuOption(choiceTransporterBack, "Back")
	utils.PrintDivider()
	fmt.Println()
}

func startTransporter() {
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
	utils.PrintInfo("IP logging is active in the ngrok logging menu")
	utils.WaitEnter()
}

func showTransporterLogs() {
	utils.PrintMenuHeader("Ngrok IP Logging")
	if !transporter.Running() {
		utils.PrintInfo("Ngrok is not running")
	}

	summaries := transporter.IPLogSummaries()
	if len(summaries) == 0 {
		utils.PrintInfo("No visits captured yet")
		utils.WaitEnter()
		return
	}

	fmt.Printf("  %-18s %-8s %-10s %-32s %s\n", "IP", "HITS", "LAST SEEN", "LAST PATH", "USER AGENT")
	utils.PrintDivider()
	for _, summary := range summaries {
		fmt.Printf(
			"  %-18s %-8d %-10s %-32s %s\n",
			summary.IP,
			summary.Hits,
			summary.LastSeen.Format("15:04:05"),
			truncate(summary.LastPath, 32),
			truncate(summary.UserAgent, 42),
		)
	}
	fmt.Println()
	utils.PrintInfo("Temporary local log only; clearing/stopping the manager does not ban anyone automatically")
	utils.WaitEnter()
}

func stopTransporter() {
	if !transporter.Running() {
		return
	}
	_ = transporter.Stop()
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	if max <= 1 {
		return value[:max]
	}
	return value[:max-3] + "..."
}
