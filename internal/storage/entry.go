package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// CreateEntry creates a new entry
func (db *DB) CreateEntry(userID int, entryType models.DataType, title, data, metadata string) (*models.Entry, error) {
	query := `
		INSERT INTO entries (user_id, type, title, data, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, type, title, data, metadata, version, updated_at, created_at
	`

	entry := &models.Entry{}
	err := db.conn.QueryRow(query, userID, entryType, title, data, metadata).Scan(
		&entry.ID, &entry.UserID, &entry.Type, &entry.Title,
		&entry.Data, &entry.Metadata, &entry.Version, &entry.UpdatedAt, &entry.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves an entry by ID and user ID
func (db *DB) GetEntry(id, userID int) (*models.Entry, error) {
	query := `
		SELECT id, user_id, type, title, data, metadata, version, updated_at, created_at
		FROM entries
		WHERE id = $1 AND user_id = $2
	`

	entry := &models.Entry{}
	err := db.conn.QueryRow(query, id, userID).Scan(
		&entry.ID, &entry.UserID, &entry.Type, &entry.Title,
		&entry.Data, &entry.Metadata, &entry.Version, &entry.UpdatedAt, &entry.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entry not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	return entry, nil
}

// ListEntries retrieves all entries for a user
func (db *DB) ListEntries(userID int) ([]models.Entry, error) {
	query := `
		SELECT id, user_id, type, title, data, metadata, version, updated_at, created_at
		FROM entries
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Type, &entry.Title,
			&entry.Data, &entry.Metadata, &entry.Version, &entry.UpdatedAt, &entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// UpdateEntry updates an existing entry
func (db *DB) UpdateEntry(id, userID int, title, data, metadata string, version int) (*models.Entry, error) {
	query := `
		UPDATE entries
		SET title = $1, data = $2, metadata = $3, version = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND user_id = $6 AND version = $7
		RETURNING id, user_id, type, title, data, metadata, version, updated_at, created_at
	`

	entry := &models.Entry{}
	err := db.conn.QueryRow(query, title, data, metadata, version+1, id, userID, version).Scan(
		&entry.ID, &entry.UserID, &entry.Type, &entry.Title,
		&entry.Data, &entry.Metadata, &entry.Version, &entry.UpdatedAt, &entry.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entry not found or version mismatch")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}

	return entry, nil
}

// DeleteEntry deletes an entry
func (db *DB) DeleteEntry(id, userID int) error {
	query := `
		DELETE FROM entries
		WHERE id = $1 AND user_id = $2
	`

	result, err := db.conn.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("entry not found")
	}

	return nil
}

// SyncEntries retrieves entries updated after a specific time
func (db *DB) SyncEntries(userID int, since time.Time) ([]models.Entry, error) {
	query := `
		SELECT id, user_id, type, title, data, metadata, version, updated_at, created_at
		FROM entries
		WHERE user_id = $1 AND updated_at > $2
		ORDER BY updated_at ASC
	`

	rows, err := db.conn.Query(query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to sync entries: %w", err)
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Type, &entry.Title,
			&entry.Data, &entry.Metadata, &entry.Version, &entry.UpdatedAt, &entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
