package service

import (
	"fmt"
	"time"

	"golang.design/x/clipboard"
)

func InitClipboard() error {
	return clipboard.Init()
}

// Copia la contraseña y  la borra despues de un minuto
func CopyPasswordToClipboard(password string) error {
	if password == "" {
		return fmt.Errorf("contraseña vacia, no se puede copiar")
	}

	clipboard.Write(clipboard.FmtText, []byte(password))

	fmt.Println("\n Contraseña copiada al portapapeles")
	fmt.Println("\n Se borrara automaticamente en 60 segundos")

	// Iniciar goroutine para borrar despues de 1 minuto

	go func() {
		time.Sleep(60 * time.Second)

		// limpiar el clipboard de forma segura
		// primero escribir datos aleatorios

		randomData := make([]byte, len(password))
		for i := range randomData {
			randomData[i] = 0
		}
		clipboard.Write(clipboard.FmtText, randomData)

		clipboard.Write(clipboard.FmtText, []byte(""))
	}()

	return nil

}

func ClearClipboard() {

	//Sobrescribir con ceros
	clipboard.Write(clipboard.FmtText, []byte{0, 0, 0, 0, 0, 0, 0, 0})

	//limpiar
	clipboard.Write(clipboard.FmtText, []byte(""))
}

// muestra cuenta regresiva en pantalla
func ShowClipboardCountdown(seconds int) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	remaining := seconds

	for range ticker.C {
		remaining -= 10

		if remaining <= 0 {
			fmt.Println("\n Portapapeles limpiado por seguridad")
			break
		}

		if remaining <= 30 {
			fmt.Printf("Portapapeles se limpiara en %d segundos...\n", remaining)
		}
	}
}
