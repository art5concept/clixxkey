package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/art5concept/clixxkey/internal/crypto"
	"github.com/art5concept/clixxkey/internal/models"
	"github.com/art5concept/clixxkey/internal/service"
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

var (
	ErrNTPFailed        = errors.New("No se pudo obtener la hora de NTP")
	ErrPasswordNotFound = errors.New("No se encontro la contraseña con ese ID")
	ErrPasswordNotReady = errors.New("La contraseña no esta lista para ser mostrada")
)

func printAllHidden(passwords []models.Password) {
	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tSitio\tUsername\tPassword")
	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")

	// el for cuando se usa con un range tiene un indice y un valor
	for _, p := range passwords {
		fmt.Fprintf(table, "%d\t%s\t%s\t****************\n", p.ID, p.Site, p.Username)
	}

	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")
	table.Flush()
}

func printTableWithUnlocked(passwords []models.Password, unlockID int, unlockedPass *models.Password) {

	table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tSitio\tUsername\tPassword")
	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")

	// el for cuando se usa con un range tiene un indice y un valor
	for _, p := range passwords {
		if p.ID == unlockID {
			// Reveal the password
			fmt.Fprintf(table, "%d\t%s\t%s\t%s\n", p.ID, p.Site, p.Username, unlockedPass.Pass)

			// limpiar memoria luego de mostrar
			passCopy := unlockedPass.Pass
			service.SecureZeroString(&passCopy)
		} else {
			// Hide the password
			fmt.Fprintf(table, "%d\t%s\t%s\t****************\n", p.ID, p.Site, p.Username)
		}
	}

	fmt.Fprintln(table, "----\t----------------\t----------------\t----------------")
	table.Flush()
}

// PrintPasswordsTable muestra las contraseñas en una tabla, desencriptando solo la que coincide con el ID y si el tiempo es correcto
// Si ID es -1, muestra todas las contraseñas ocultas

func PrintPasswordsTable(passwords []models.Password, id int) (Unlocked bool, err error) {

	if id == -1 {
		printAllHidden(passwords)
		return false, nil
	}

	status, err := service.CheckPasswordUnlock(passwords, id)

	if status.CanUnlock {
		// mostrar contraseña revelada
		printTableWithUnlocked(passwords, id, status.Password)

		if err := service.CopyPasswordToClipboard(status.Password.Pass); err != nil {
			fmt.Printf("Advertencia: No se pudo copiar al portapapeles: %v\n", err)
		}

	} else {
		printAllHidden(passwords)
	}

	// mostrar el tiempo restante
	service.PrintUnlockTimeInfo(status)

	return status.CanUnlock, err
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
