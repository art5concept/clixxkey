package service

import (
	"fmt"
	"os"

	"bufio"

	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"golang.org/x/term"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func GetNTPTime() (time.Time, error) {
	return ntp.Time("pool.ntp.org")
}

func UpdateUnlockTime() (unlockTime time.Time, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Periodo de bloqueo (ejemplo: 10s, 30d, 2m, 3y): ")
	scanner.Scan()
	period := strings.TrimSpace(scanner.Text())

	ntpTime, err := GetNTPTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("error obteniendo hora NTP: %w", err)
	}

	if len(period) < 2 {
		return time.Time{}, fmt.Errorf("formato de periodo inválido: %s", period)
	}

	value, err := strconv.Atoi(period[:len(period)-1])
	if err != nil {
		return time.Time{}, fmt.Errorf("error convirtiendo valor: %w", err)
	}

	// operador slicing o rebanado
	// desde:hasta

	unit := period[len(period)-1:]

	// Verifica el sufijo y convierte el número

	switch unit {
	case "s":
		return ntpTime.Add(time.Duration(value) * time.Second), nil
	case "d":
		return ntpTime.AddDate(0, 0, value), nil
	case "m":
		return ntpTime.AddDate(0, value, 0), nil
	case "y":
		return ntpTime.AddDate(value, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unidad inválida: %s, usa: s, d, m, y", unit)
	}
}

func ReadPassword() ([]byte, error) {
	// fmt.Print("Introduce tu contraseña maestra: ")
	termState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	defer term.Restore(int(os.Stdin.Fd()), termState)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	// fmt.Println()
	if err != nil {
		return nil, err
	}

	return password, nil
}

func SecureZeroString(s *string) {
	if s == nil {
		return
	}
	b := []byte(*s)
	for i := range b {
		b[i] = 0
	}
	*s = string(b)
}
