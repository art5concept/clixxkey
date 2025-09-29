//go:build !windows

package service

// En sistemas no Windows, la función es vacía
func EnableVirtualTerminalProcessing() {}
