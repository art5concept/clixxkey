//go:build windows

package service

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func EnableVirtualTerminalProcessing() {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	stdout := os.Stdout.Fd()
	var mode uint32
	getConsoleMode.Call(uintptr(stdout), uintptr(unsafe.Pointer(&mode)))
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	setConsoleMode.Call(uintptr(stdout), uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}
