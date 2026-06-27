package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)


func PrintBanner(){
	banner:=
`
  ______    ______          ______    ______   _______   ______  __    __   ______  
 /      \  /      \        /      \  /      \ |       \ |      \|  \  |  \ /      \ 
|  $$$$$$\|  $$$$$$\      |  $$$$$$\|  $$$$$$\| $$$$$$$\ \$$$$$$| $$\ | $$|  $$$$$$\
| $$___\$$| $$___\$$      | $$   \$$| $$  | $$| $$  | $$  | $$  | $$$\| $$| $$ __\$$
 \$$    \  \$$    \       | $$      | $$  | $$| $$  | $$  | $$  | $$$$\ $$| $$|    \
 _\$$$$$$\ _\$$$$$$\      | $$   __ | $$  | $$| $$  | $$  | $$  | $$\$$ $$| $$ \$$$$
|  \__| $$|  \__| $$      | $$__/  \| $$__/ $$| $$__/ $$ _| $$_ | $$ \$$$$| $$__| $$
 \$$    $$ \$$    $$       \$$    $$ \$$    $$| $$    $$|   $$ \| $$  \$$$ \$$    $$
  \$$$$$$   \$$$$$$         \$$$$$$   \$$$$$$  \$$$$$$$  \$$$$$$ \$$   \$$  \$$$$$$ 
                                                                                    
                                                                                    

`
	fmt.Print(HiGold.Apply(banner))
}

func ClearTerminal(){
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin", "linux":
		fmt.Print("\033[?1049l")
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

// EnterAltScreen switches to the terminal alternate buffer so live bot
// output does not overwrite the main menu screen.
func EnterAltScreen() {
	if runtime.GOOS == "windows" {
		ClearTerminal()
		return
	}
	fmt.Print("\033[?1049h\033[H\033[2J")
}

// LeaveAltScreen restores the main terminal buffer.
func LeaveAltScreen() {
	if runtime.GOOS == "windows" {
		return
	}
	fmt.Print("\033[?1049l")
}

