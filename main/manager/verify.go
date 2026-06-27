package main

import (
	"fmt"

	"ss-coding/utils"
)

func runVerify() {
	dir, err := utils.WebAppDir()
	if err != nil {
		utils.PrintError(err.Error())
		utils.WaitEnter()
		return
	}

	utils.PrintInfo(fmt.Sprintf("Scanning %s", dir))
	fmt.Println()

	checks := []struct {
		label string
		args  []string
	}{
		{"TypeScript type check", []string{"run", "typecheck"}},
		{"ESLint", []string{"run", "lint"}},
		{"Dependency audit", []string{"audit"}},
	}

	passed := 0
	for _, check := range checks {
		result := utils.RunCapture(dir, "bun", check.args...)
		if utils.PrintStep(check.label, result) {
			passed++
		}
		fmt.Println()
	}

	summary := fmt.Sprintf("%d/%d checks passed", passed, len(checks))
	if passed == len(checks) {
		utils.PrintSuccess(summary)
	} else {
		utils.PrintError(summary)
	}
	utils.WaitEnter()
}
