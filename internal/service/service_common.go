package service

import (
	"fmt"
	"os"

	"bufio"

	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
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
	period := scanner.Text()

	ntpTime, err := GetNTPTime()
	if err != nil {
		fmt.Println("No se pudo obtener la hora NTP:", err)
		// continue
	}

	if strings.HasSuffix(period, "s") {
		secs, _ := strconv.Atoi(strings.TrimSuffix(period, "s"))
		return ntpTime.Add(time.Duration(secs) * time.Second), nil
	} else if strings.HasSuffix(period, "d") {
		days, _ := strconv.Atoi(strings.TrimSuffix(period, "d"))
		return ntpTime.AddDate(0, 0, days), nil
	} else if strings.HasSuffix(period, "m") {
		month, _ := strconv.Atoi(strings.TrimSuffix(period, "m"))
		return ntpTime.AddDate(0, month, 0), nil
	} else if strings.HasSuffix(period, "y") {
		years, _ := strconv.Atoi(strings.TrimSuffix(period, "y"))
		return ntpTime.AddDate(years, 0, 0), nil
	} else {
		fmt.Println("Formato de periodo inválido.")
		return time.Time{}, err
	}
}
