package models

import "time"

// DataType represents the type of data being stored
type DataType string

const (
	DataTypeLogin  DataType = "login"
	DataTypeText   DataType = "text"
	DataTypeBinary DataType = "binary"
	DataTypeCard   DataType = "card"
)

// Entry represents a single data entry in the system
type Entry struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      DataType  `json:"type"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`     // Base64-encoded encrypted data
	Metadata  string    `json:"metadata"` // JSON metadata
	Version   int       `json:"version"`  // For conflict resolution
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginEntry represents login/password data
type LoginEntry struct {
	Site     string `json:"site"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Notes    string `json:"notes,omitempty"`
}

// TextEntry represents arbitrary text data
type TextEntry struct {
	Content string `json:"content"`
	Notes   string `json:"notes,omitempty"`
}

// BinaryEntry represents binary data
type BinaryEntry struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	Notes    string `json:"notes,omitempty"`
}

// CardEntry represents bank card information
type CardEntry struct {
	Number      string `json:"number"`
	Holder      string `json:"holder"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	Notes       string `json:"notes,omitempty"`
}

// CreateEntryRequest represents a request to create a new entry
type CreateEntryRequest struct {
	Type     DataType `json:"type"`
	Title    string   `json:"title"`
	Data     string   `json:"data"`
	Metadata string   `json:"metadata"`
}

// UpdateEntryRequest represents a request to update an existing entry
type UpdateEntryRequest struct {
	Title    string `json:"title"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	Version  int    `json:"version"`
}

// SyncResponse represents the result of a sync operation
type SyncResponse struct {
	Entries []Entry `json:"entries"`
	Version int     `json:"version"`
}
