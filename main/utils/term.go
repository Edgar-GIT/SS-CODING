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
	fmt.Print(Yellow.Apply(banner))
}

func ClearTerminal(){
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin", "linux":
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
