package file

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/art5concept/clixxkey/internal/crypto"
	"github.com/art5concept/clixxkey/internal/models"
)

type FileRepository struct {
	path string
	key  []byte
}

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

func (r FileRepository) List() ([]models.Password, error) {
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

func PrintPasswordsTable(passwords []models.Password, id int) (Unlocked bool) {
	// scanner := bufio.NewScanner(os.Stdin)
	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tSitio\tUsername\tPassword")
	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")

	// el for cuando se usa con un range tiene un indice y un valor
	for _, p := range passwords {
		if p.ID == id && time.Now().After(p.UnlockAfter) {
			fmt.Fprintf(table, "%d\t%s\t%s\t%s\n", p.ID, p.Site, p.Username, p.Pass)
			Unlocked = true
		} else {
			fmt.Fprintf(table, "%d\t%s\t%s\t****************\n", p.ID, p.Site, p.Username)
			Unlocked = false
		}
		// %d es para enteros, %s para strings
		// \t es tabulador
		// \n es nueva linea

	}
	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")
	table.Flush()
	fmt.Println("\tListados exitosamente")
	return Unlocked
}

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

}

func (fr *FileRepository) UpdateUnlockAfter(id int, unlockAfter time.Time) error {
	passwords, err := fr.List()
	if err != nil {
		return err
	}

	for i, p := range passwords {
		if p.ID == id {
			passwords[i].UnlockAfter = unlockAfter
			break
		}
	}

	data, err := json.MarshalIndent(passwords, "", "  ")
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(data, fr.key)
	if err != nil {
		return err
	}
	return os.WriteFile(fr.path, encrypted, 0600)
}
