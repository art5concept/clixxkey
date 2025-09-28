package file

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/art5concept/clixxkey/internal/crypto"
	"github.com/art5concept/clixxkey/internal/models"
)

type FileRepository struct {
	path string
	key  []byte
}

// type EncryptedFileRepository interface {
// 	NewEncrypted(path string, key []byte) *FileRepository
// 	Delete(id int) error
// }

func New(path string) *FileRepository {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}
	return &FileRepository{path: path}
}

func NewEncrypted(path string, key []byte) *FileRepository {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte{}, 0600)
	}
	return &FileRepository{path: path, key: key}
}

// func (fr *FileRepository) List() ([]models.Password, error) {
// 	data, err := os.ReadFile(fr.path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var passwords []models.Password
// 	err = json.Unmarshal(data, &passwords)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return passwords, nil
// }

func (r *FileRepository) List() ([]models.Password, error) {
	encrypted, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	if len(encrypted) == 0 {
		return []models.Password{}, nil
	}
	data, err := crypto.Decrypt(encrypted, r.key)
	if err != nil {
		return nil, err
	}
	var passwords []models.Password
	if err := json.Unmarshal(data, &passwords); err != nil {
		return nil, err
	}
	return passwords, nil
}

func PrintPasswordsTable(passwords []models.Password) {
	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tSitio\tUsername")
	fmt.Fprintln(table, "----\t----------------\t----------------")
	for _, p := range passwords {
		fmt.Fprintf(table, "%d\t%s\t%s\n", p.ID, p.Site, p.Username)
	}
	table.Flush()
	fmt.Println("----------------------------------------")
	fmt.Println("\tListados exitosamente")
}

// func (fr *FileRepository) Save(p models.Password) error {
// 	passwords, err := fr.List()
// 	if err != nil {
// 		return err
// 	}

// 	// Assign a new ID
// 	var maxID int
// 	for _, pass := range passwords {
// 		if pass.ID > maxID {
// 			maxID = pass.ID
// 		}
// 	}

// 	p.ID = maxID + 1

// 	passwords = append(passwords, p)
// 	data, err := json.MarshalIndent(passwords, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(fr.path, data, 0644)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *FileRepository) Save(p models.Password) error {
	passwords, err := r.List()
	if err != nil {
		return err
	}

	// Assign a new ID
	var maxID int
	for _, pass := range passwords {
		if pass.ID > maxID {
			maxID = pass.ID
		}
	}

	p.ID = maxID + 1

	passwords = append(passwords, p)
	data, err := json.Marshal(passwords)
	if err != nil {
		return err
	}
	encrypted, err := crypto.Encrypt(data, r.key)
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, encrypted, 0600)
}

func (fr *FileRepository) Delete(id int) error {
	passwords, err := fr.List()
	if err != nil {
		return err
	}

	var updatedPasswords []models.Password
	for _, p := range passwords {
		if p.ID != id {
			updatedPasswords = append(updatedPasswords, p)
		}
	}

	data, err := json.MarshalIndent(updatedPasswords, "", "  ")
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(data, fr.key)
	if err != nil {
		return err
	}
	return os.WriteFile(fr.path, encrypted, 0600)

	// err = os.WriteFile(fr.path, data, 0644)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

// CRUD operations for file-based repository would go here

// func OpenFile(filename string) (*os.File, error) {
// 	// file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
// 	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return nil, err
// 	}
// 	return file, nil
// }

// func ManualCloseFile(file *os.File) {
// 	err := file.Close()
// 	if err != nil {
// 		fmt.Println("Error closing file:", err)
// 		return
// 	}
// 	fmt.Println("File closed successfully")
// }

// // ReadFile reads and decodes passwords from the given file

// func DecodePasswordsFromFile(file *os.File) ([]models.Password, error) {
// 	var passwords []models.Password
// 	decoder := json.NewDecoder(file)
// 	err := decoder.Decode(&passwords)
// 	if err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return nil, err
// 	}
// 	return passwords, nil
// }

// func PrintFile() {

// 	pass1 := models.Password{
// 		ID:       1,
// 		Site:     "example.com",
// 		Username: "user1",
// 		Pass:     "pass1",
// 	}

// 	passenc, err := json.Marshal(pass1)
// 	if err != nil {
// 		fmt.Println("Error marshalling password:", err)
// 	}

// 	fmt.Println(string(passenc))
// }

// func CreateFile() {
// 	file, err := os.Create("passwords.json")
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer file.Close()
// 	fmt.Println("File created successfully")
// }

// func UpdateFile() {
// 	file, err := os.Open("passwords.json")
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	var passwords []models.Password
// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&passwords)
// 	if err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		return
// 	}

// 	// Example update: change password for ID 1
// 	for i, p := range passwords {
// 		if p.ID == 1 {
// 			passwords[i].Pass = "newpass1"
// 		}
// 	}

// 	file.Close() // Close the file before reopening for writing

// 	file, err = os.Create("passwords.json")

// 	if err != nil {
// 		fmt.Println("Error opening file for writing:", err)
// 		return
// 	}
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(passwords)
// 	if err != nil {
// 		fmt.Println("Error encoding JSON:", err)
// 		return
// 	}

// 	fmt.Println("Password updated successfully")
// }

// func DeleteFile() {
// 	// Implementation for deleting a password entry from the file
// }
