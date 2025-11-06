package models

import (
	"testing"
)

func TestDataTypeString(t *testing.T) {
	tests := []struct {
		dt       DataType
		expected string
	}{
		{DataTypeLogin, "login"},
		{DataTypeText, "text"},
		{DataTypeBinary, "binary"},
		{DataTypeCard, "card"},
	}

	for _, tt := range tests {
		if string(tt.dt) != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, string(tt.dt))
		}
	}
}

func TestEntryStructure(t *testing.T) {
	entry := Entry{
		ID:       1,
		UserID:   123,
		Type:     DataTypeLogin,
		Title:    "Test Entry",
		Data:     "encrypted_data",
		Metadata: "{\"site\":\"example.com\"}",
		Version:  1,
	}

	if entry.ID != 1 {
		t.Errorf("Expected ID 1, got %d", entry.ID)
	}

	if entry.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", entry.UserID)
	}

	if entry.Type != DataTypeLogin {
		t.Errorf("Expected Type login, got %s", entry.Type)
	}
}

func TestLoginEntryStructure(t *testing.T) {
	login := LoginEntry{
		Site:     "example.com",
		Login:    "user@example.com",
		Password: "secret123",
		Notes:    "Personal account",
	}

	if login.Site != "example.com" {
		t.Errorf("Expected site example.com, got %s", login.Site)
	}
}
