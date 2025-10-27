package storage

import (
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// Storage defines the interface for data storage operations
// This allows us to use mock implementations in tests
type Storage interface {
	// User operations
	CreateUser(username, passwordHash string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	UserExists(username string) (bool, error)

	// Entry operations
	CreateEntry(userID int, entryType models.DataType, title, data, metadata string) (*models.Entry, error)
	GetEntry(id, userID int) (*models.Entry, error)
	ListEntries(userID int) ([]models.Entry, error)
	UpdateEntry(id, userID int, title, data, metadata string, version int) (*models.Entry, error)
	DeleteEntry(id, userID int) error
	SyncEntries(userID int, since time.Time) ([]models.Entry, error)
}

// Ensure DB implements Storage interface
var _ Storage = (*DB)(nil)
