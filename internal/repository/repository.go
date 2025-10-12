package repository

import (
	"time"

	"github.com/art5concept/clixxkey/internal/models"
)

type Repository interface {
	List() ([]models.Password, error)
	Save(p models.Password) error
	Delete(id int) error
	UpdateUnlockAfter(id int, unlockAfter time.Time) error

	// Update(id int, p models.Password) error
	// GetByID(id int) (models.Password, error)
}

// in the future maybe ill make a GUI to use the interface here and also i want to put some reed-salomon error manager here
