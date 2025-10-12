package service

import (
	"fmt"
)

func EnterAltScreen() {
	fmt.Print("\033[?1049h")
	fmt.Print("\x1b[H")
}

// I use AI in specially this parts of the code because i dont know how do this exit stuff
func ExitAltScreen() {
	fmt.Print("\x1b[H")
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[3J")
	fmt.Print("\x1b[?1049l")
}
