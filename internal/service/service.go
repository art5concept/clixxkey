//go:build !windows

package service

import "fmt"

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// En sistemas no Windows, la función es vacía
func EnableVirtualTerminalProcessing() {}
