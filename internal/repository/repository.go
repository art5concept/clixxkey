package repository

import "github.com/art5concept/clixxkey/internal/models"

type Repository interface {
	List() ([]models.Password, error)
	Save(p models.Password) error
	Delete(id int) error

	// Update(id int, p models.Password) error
	// GetByID(id int) (models.Password, error)
}
