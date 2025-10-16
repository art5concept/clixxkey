package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/art5concept/clixxkey/internal/crypto"
	"github.com/art5concept/clixxkey/internal/models"
	"github.com/art5concept/clixxkey/internal/repository"
	"github.com/art5concept/clixxkey/internal/repository/file"
	"github.com/art5concept/clixxkey/internal/service"
)

var repo repository.Repository

func main() {

	service.EnableVirtualTerminalProcessing()

	service.EnterAltScreen()
	defer service.ExitAltScreen()

	fmt.Println("CLI password manager")

	fmt.Print("Introduce tu contraseña maestra: ")
	scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// password := scanner.Text()

	password, err := service.ReadPassword()
	if err != nil {
		fmt.Println("Error leyendo contraseña:", err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("No se pudo obtener el home:", err)
		os.Exit(1)
	}
	dataDir := home + "/.clixxkey"
	os.MkdirAll(dataDir, 0700)

	jsonPath := dataDir + "/passwords.json"
	saltPath := dataDir + "/passwords.salt"

	// Salt: genera o lee de archivo
	var salt []byte
	if _, err := os.Stat(saltPath); os.IsNotExist(err) {
		salt, _ = crypto.GenerateSalt()
		os.WriteFile(saltPath, salt, 0600)
	} else {
		salt, err = os.ReadFile(saltPath)
		if err != nil {
			fmt.Println("Error leyendo salt:", err)
			os.Exit(1)
		}
	}

	key := crypto.DeriveKey([]byte(password), salt)
	// service.SecureZeroString(&password)
	repo = file.NewEncrypted(jsonPath, key)

	_, err = repo.List()
	if err != nil {
		fmt.Println("Error de desencriptado.")
		time.Sleep(2 * time.Second)
		os.Exit(1)
	}

	scanner = bufio.NewScanner(os.Stdin)

	for {

		// this is my menu all made by me and to me looks simple and clean
		service.ClearScreen()
		fmt.Println("----------------------------------------")
		fmt.Println("\n Password Manager ")
		fmt.Println("1. Mostrar Contraseñas")
		fmt.Println("2. Agregar Contraseña")
		fmt.Println("3. Borrar Contraseña")
		fmt.Println("0. Salirr")
		fmt.Println("Seleciona una opcion:")

		scanner.Scan()
		option := scanner.Text()

		if option == "0" {
			fmt.Println("Exiting...")
			break
		} else if option == "1" {
			passwords, err := repo.List()
			if err != nil {
				fmt.Println("Error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			service.ClearScreen()
			// Mostrar tabla de contraseñas sin revelar contraseñas
			file.PrintPasswordsTable(passwords, -1)

			fmt.Println("\nIngrese el Id para ver el secreto o presione Enter para regresar:")
			scanner.Scan()
			idStr := scanner.Text()

			// si presiona enter sin escribir nada regresa al menu
			if idStr == "" {
				continue
			}

			idInt, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Id invalido, debe ser un numero.")
				time.Sleep(2 * time.Second)
				continue
			}

			service.ClearScreen()

			// intenta desbloquear id especifico
			unlocked, err := file.PrintPasswordsTable(passwords, idInt)

			// time.Sleep(2 * time.Second)

			if err != nil && !errors.Is(err, service.ErrPasswordLocked) {
				//Errores criticos
				time.Sleep(2 * time.Second)
				continue
			}

			// si llegamos aqui, la contraseña  unlocked == true y se mostro
			// time.Sleep(2 * time.Second)

			if unlocked {
				fmt.Println("\n La contraseña se mostrara por unos segundos...")
				time.Sleep(5 * time.Second)

				service.ClearScreen()

				unlockAfter, err := service.UpdateUnlockTime()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error actualizando tiempo de desbloqueo: %v\n", err)
					time.Sleep(2 * time.Second)
					continue
				}

				err = repo.UpdateUnlockAfter(idInt, unlockAfter)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error guardando nuevo tiempo de desbloqueo: %v\n", err)
					time.Sleep(2 * time.Second)
					continue
				}
				fmt.Println("-----------------------------------------------")
				fmt.Println("\n El tiempo de desbloqueo se ha actualizado exitosamente.")
				fmt.Println(" La contraseña permanecera en el portapapeles por 60 segundos")
				fmt.Println(" Press Enter to return to main menu.")
				scanner.Scan()
			} else {
				time.Sleep(3 * time.Second)
			}
			continue

		} else if option == "2" {
			service.ClearScreen()

			fmt.Print("Sitio web: ")
			scanner.Scan()
			site := scanner.Text()

			fmt.Print("Nombre de usuario: ")
			scanner.Scan()
			username := scanner.Text()

			fmt.Print("Contraseña: ")
			scanner.Scan()
			pass := scanner.Text()

			unlockAfter, err := service.UpdateUnlockTime()

			if err != nil {
				fmt.Println("Error updating unlock time:", err)
				continue
			}

			newPass := models.Password{
				ID:          0, // ID will be set in Save method
				Site:        site,
				Username:    username,
				Pass:        pass,
				UnlockAfter: unlockAfter,
			}

			if err := repo.Save(newPass); err != nil {
				fmt.Println("Error de guardado:", err)
				continue
			} else {
				fmt.Println("----------------------------------------")
				fmt.Println("Contraseña guardada exitosamente.")
				continue
			}

		} else if option == "3" {
			passwords, err := repo.List()
			if err != nil {
				fmt.Println("Error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			service.ClearScreen()
			// Mostrar tabla de contraseñas sin revelar contraseñas
			file.PrintPasswordsTable(passwords, -1)

			fmt.Print("Ingrese el ID de la contraseña a borrar o Enter para regresar: ")
			scanner.Scan()
			idString, err := strconv.Atoi(scanner.Text())

			if err != nil {
				fmt.Println("Invalid ID, please enter a number.")
				continue
			}

			if err := repo.Delete(idString); err != nil {
				fmt.Println("Error al borrar:", err)
				continue
			} else {
				fmt.Println("Contraseña borrada exitosamente.")
				continue
			}

		} else {
			fmt.Println("Invalid option, please try again.")
			continue
		}

	}

}

// i want too say thanks to all the CS50 staff and harvard univerity
