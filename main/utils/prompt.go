package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintDivider() {
	fmt.Println(HiBlack.Apply(strings.Repeat("─", 52)))
}

func PrintMenuHeader(title string) {
	PrintDivider()
	fmt.Println(BoldHiCyan.Apply("  " + title))
	PrintDivider()
}

func PrintMenuOption(key, label string) {
	fmt.Printf("  %s  %s\n", BoldHiYellow.Apply(key+"."), label)
}

func ReadChoice(prompt string) string {
	fmt.Print(BoldHiWhite.Apply(prompt))
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func PrintSuccess(message string) {
	fmt.Println(HiGreen.Apply("✓ " + message))
}

func PrintError(message string) {
	fmt.Println(HiRed.Apply("✗ " + message))
}

func PrintInfo(message string) {
	fmt.Println(HiCyan.Apply("→ " + message))
}

func WaitEnter() {
	fmt.Print(Muted("Press Enter to continue..."))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func WaitForStop() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
