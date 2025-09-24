package main

import (
	"fmt"
	"os"

	"github.com/art5concept/clixxkey/internal/repository/file"
)

func main() {

	fmt.Println("CLI password manager")
	// file.PrintFile()

	// file.CreateFile()

	f, err := file.OpenFile("passwords.json")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// defer file.ManualCloseFile(f)
	defer f.Close()

	exist, err := os.Stat("passwords.json")

	if err != nil || exist.Size() == 0 {
		fmt.Println("File does not exist or is empty, creating a new one...")
		// Create an empty JSON array in the file
		_, err := f.WriteString("[]")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		// Reset file pointer to the beginning after writing
		_, err = f.Seek(0, 0)
		if err != nil {
			fmt.Println("Error seeking file:", err)
			return
		}
	}

	passwords, err := file.DecodePasswordsFromFile(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	file.PrintPasswordsTable(passwords)

}
