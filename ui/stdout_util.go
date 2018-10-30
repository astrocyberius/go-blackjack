package ui

import (
	"os"
	"os/exec"
	"runtime"
)

// Create a map for storing clear functions
var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["darwin"] = clear["linux"]

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// ClearScreen clears the Stdout screen
func ClearScreen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear the terminal screen :(")
	}
}
