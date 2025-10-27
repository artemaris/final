package storage

import (
	"database/sql"
	"fmt"

	"github.com/artemaris/gophkeeper/internal/models"
)

// CreateUser creates a new user
func (db *DB) CreateUser(username, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, username, created_at
	`

	user := &models.User{}
	err := db.conn.QueryRow(query, username, passwordHash).Scan(
		&user.ID, &user.Username, &user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := db.conn.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := db.conn.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UserExists checks if a user exists
func (db *DB) UserExists(username string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
	`
	var exists bool
	err := db.conn.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}
