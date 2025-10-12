package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/art5concept/clixxkey/internal/models"
)

var (
	ErrNTPFailed        = errors.New("no se pudo obtener la hora de NTP")
	ErrPasswordNotFound = errors.New("el ID de contraseña no existe")
	ErrPasswordLocked   = errors.New("la contraseña aún está bloqueada")
)

type PasswordUnlockStatus struct {
	CanUnlock     bool
	Password      *models.Password
	TimeRemaining time.Duration
	Error         error
	CurrentTime   time.Time
}

func CheckPasswordUnlock(passwords []models.Password, id int) (*PasswordUnlockStatus, error) {
	status := &PasswordUnlockStatus{
		CanUnlock: false,
	}

	var targetPassword *models.Password

	// Buscar la contraseña con el ID dado
	for i := range passwords {
		if passwords[i].ID == id {
			targetPassword = &passwords[i]
			break
		}
	}

	if targetPassword == nil {
		status.Error = ErrPasswordNotFound
		return status, ErrPasswordNotFound
	}

	status.Password = targetPassword

	// Obtener la hora actual desde un servidor NTP
	realTime, err := GetNTPTime()
	if err != nil {
		status.Error = fmt.Errorf("%w: %v", ErrNTPFailed, err)
		return status, status.Error
	}
	status.CurrentTime = realTime

	// Comparar la hora actual con UnlockAfter
	if realTime.After(targetPassword.UnlockAfter) {
		status.CanUnlock = true
		status.TimeRemaining = 0
		return status, nil
	}

	status.TimeRemaining = targetPassword.UnlockAfter.Sub(realTime)
	status.Error = ErrPasswordLocked
	return status, ErrPasswordLocked
}

// this function was also almost made by me in fully
func FormatTimeRemaining(d time.Duration) string {
	if d <= 0 {
		return "ya desbloqueada"
	}

	totalSeconds := int(d.Seconds())

	years := totalSeconds / (365 * 24 * 3600)
	totalSeconds %= 365 * 24 * 3600

	months := totalSeconds / (30 * 24 * 3600)
	totalSeconds %= 30 * 24 * 3600

	days := totalSeconds / (24 * 3600)
	totalSeconds %= 24 * 3600

	hours := totalSeconds / 3600
	totalSeconds %= 3600

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d años", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d meses", months))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d días", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d horas", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutos", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d segundos", seconds))
	}

	if len(parts) == 1 {
		return parts[0]
	} else if len(parts) == 2 {
		return parts[0] + " y " + parts[1]
	} else {
		return parts[0] + " y " + parts[1]
	}
}

func PrintUnlockTimeInfo(status *PasswordUnlockStatus) {
	if status.CanUnlock {
		fmt.Println("\n Contraseña disponible para desbloquear.")
		return
	}

	if status.Error != nil {
		if errors.Is(status.Error, ErrPasswordLocked) {
			fmt.Printf("\n Tiempo restante para desbloquear: %s\n", FormatTimeRemaining(status.TimeRemaining))
			fmt.Printf(" Se desbloqueara el: %s\n", status.Password.UnlockAfter.Format("2006-01-02 15:04:05 MST"))
		} else if errors.Is(status.Error, ErrNTPFailed) {
			fmt.Println("\n Error de sincronizacion de tiempo")
			fmt.Println("No se puede verificar el tiempo de desbloqueo por seguridad")
		} else if errors.Is(status.Error, ErrPasswordNotFound) {
			fmt.Println("\n Contraseña no encontrada")
		}
	}
}
