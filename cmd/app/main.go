package main

import (
	"bufio"
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

	fmt.Println("CLI password manager")

	fmt.Print("Introduce tu contraseña maestra: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	password := scanner.Text()

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
		salt, _ = os.ReadFile(saltPath)
	}

	key := crypto.DeriveKey([]byte(password), salt)

	repo = file.NewEncrypted(jsonPath, key)

	_, err = repo.List()
	if err != nil {
		fmt.Println("Error de desencriptado.")
		time.Sleep(2 * time.Second)
		os.Exit(1)
	}

	scanner = bufio.NewScanner(os.Stdin)

	for {
		service.ClearScreen()
		fmt.Println("----------------------------------------")
		fmt.Println("\n Password Manager ")
		fmt.Println("1. Mostrar Contraseñas")
		fmt.Println("2. Agregar Contraseña")
		fmt.Println("3. Borrar Contraseña")
		fmt.Println("0. Salir")
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
				return
			}
			service.ClearScreen()
			file.PrintPasswordsTable(passwords, -1)
			fmt.Println("Enter ID to view password details or press Enter to return:")
			scanner.Scan()
			id := scanner.Text()
			idInt, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("Invalid ID, please enter a number.")
				continue
			}
			service.ClearScreen()
			Unlocked := file.PrintPasswordsTable(passwords, idInt)
			time.Sleep(2 * time.Second)

			if Unlocked {
				unlockAfter, err := service.UpdateUnlockTime()
				if err != nil {
					fmt.Println("Error updating unlock time:", err)
					continue
				}

				repo.UpdateUnlockAfter(idInt, unlockAfter)

				fmt.Println("Press Enter to return to main menu.")
				scanner.Scan()

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
			service.ClearScreen()
			fmt.Print("Ingrese el ID de la contraseña a borrar: ")
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
