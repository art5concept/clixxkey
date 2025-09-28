package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

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
	// file.PrintFile()

	// file.CreateFile()

	fmt.Print("Introduce tu contraseña maestra: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	password := scanner.Text()

	// Salt: genera o lee de archivo
	var salt []byte
	if _, err := os.Stat("passwords.salt"); os.IsNotExist(err) {
		salt, _ = crypto.GenerateSalt()
		os.WriteFile("passwords.salt", salt, 0600)
	} else {
		salt, _ = os.ReadFile("passwords.salt")
	}

	key := crypto.DeriveKey([]byte(password), salt)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("No se pudo obtener el home:", err)
		os.Exit(1)
	}
	dataDir := home + "/.clixxkey"
	os.MkdirAll(dataDir, 0700) // Crea la carpeta con permisos seguros

	jsonPath := dataDir + "/passwords.json"
	saltPath := dataDir + "/passwords.salt"

	repo = file.NewEncrypted(jsonPath, key)
	scanner = bufio.NewScanner(os.Stdin)

	for {
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
			file.PrintPasswordsTable(passwords)
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

			newPass := models.Password{
				ID:       0, // ID will be set in Save method
				Site:     site,
				Username: username,
				Pass:     pass,
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

		//

		// if !scanner.Scan() {
		// 	break
		// }
	}

	// pass := models.Password{
	// 	ID:       0,
	// 	Site:     "example.com",
	// 	Username: "user123",
	// 	Pass:     "securepassword",
	// }

	// repo.Save(pass)

	// passwords, err := repo.List()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// file.PrintPasswordsTable(passwords)

	// Asegura que "repo" implementa el Repository interface
	var _ repository.Repository = repo

	// passwords, err := repo.List()

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// // defer file.ManualCloseFile(f)
	// defer f.Close()

	// exist, err := os.Stat("passwords.json")

	// if err != nil || exist.Size() == 0 {
	// 	fmt.Println("File does not exist or is empty, creating a new one...")
	// 	// Create an empty JSON array in the file
	// 	_, err := f.WriteString("[]")
	// 	if err != nil {
	// 		fmt.Println("Error writing to file:", err)
	// 		return
	// 	}
	// 	// Reset file pointer to the beginning after writing
	// 	_, err = f.Seek(0, 0)
	// 	if err != nil {
	// 		fmt.Println("Error seeking file:", err)
	// 		return
	// 	}
	// }

	// passwords, err := file.DecodePasswordsFromFile(f)
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// file.PrintPasswordsTable(passwords)

}
