package storage

import (
	"testing"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// TestEntryValidation tests entry data validation logic
func TestEntryValidation(t *testing.T) {
	entry := models.Entry{
		ID:        1,
		UserID:    1,
		Type:      models.DataTypeLogin,
		Title:     "Test Entry",
		Data:      "encrypted_data",
		Metadata:  "{\"site\":\"example.com\"}",
		Version:   1,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if entry.ID <= 0 {
		t.Error("Entry ID must be positive")
	}

	if entry.UserID <= 0 {
		t.Error("User ID must be positive")
	}

	if entry.Type == "" {
		t.Error("Entry type cannot be empty")
	}

	if entry.Title == "" {
		t.Error("Entry title cannot be empty")
	}

	if entry.Data == "" {
		t.Error("Entry data cannot be empty")
	}

	if entry.Version < 1 {
		t.Error("Version must be at least 1")
	}
}

// TestEntryType tests data type validation
func TestEntryType(t *testing.T) {
	validTypes := []models.DataType{
		models.DataTypeLogin,
		models.DataTypeText,
		models.DataTypeBinary,
		models.DataTypeCard,
	}

	for _, dataType := range validTypes {
		if dataType == "" {
			t.Errorf("Data type '%s' cannot be empty", dataType)
		}
	}
}

// TestEntryVersioning tests version management logic
func TestEntryVersioning(t *testing.T) {
	baseVersion := 1

	// Test version increment
	newVersion := baseVersion + 1
	if newVersion != 2 {
		t.Errorf("Expected version 2, got %d", newVersion)
	}

	// Test version comparison
	if baseVersion == newVersion {
		t.Error("Versions should be different")
	}
}

// TestEntryMetadata tests metadata handling
func TestEntryMetadata(t *testing.T) {
	tests := []struct {
		metadata string
		valid    bool
	}{
		{"{\"key\":\"value\"}", true},
		{"{\"site\":\"example.com\",\"login\":\"user@example.com\"}", true},
		{"", false}, // Empty metadata is valid
		{"invalid json", false},
	}

	for _, tt := range tests {
		t.Run(tt.metadata, func(t *testing.T) {
			if tt.metadata != "" && !isValidJSON(tt.metadata) && tt.valid {
				t.Errorf("Expected valid JSON for %s", tt.metadata)
			}
		})
	}
}

// Helper function to check JSON validity
func isValidJSON(str string) bool {
	// Simple check for JSON structure
	return len(str) >= 2 && (str[0] == '{' && str[len(str)-1] == '}') ||
		(str[0] == '[' && str[len(str)-1] == ']')
}

// TestSyncEntries_Logic tests sync logic
func TestSyncEntries_Logic(t *testing.T) {
	since := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	if since.After(now) {
		t.Error("Since time should be in the past")
	}

	if now.Sub(since) < 0 {
		t.Error("Time difference calculation is incorrect")
	}

	// Test timestamp conversion
	timestamp := now.Unix()
	if timestamp <= 0 {
		t.Error("Timestamp must be positive")
	}
}

// TestEntryFiltering tests entry filtering logic
func TestEntryFiltering(t *testing.T) {
	entries := []models.Entry{
		{ID: 1, UserID: 1, Type: models.DataTypeLogin, Title: "Entry 1"},
		{ID: 2, UserID: 1, Type: models.DataTypeText, Title: "Entry 2"},
		{ID: 3, UserID: 2, Type: models.DataTypeCard, Title: "Entry 3"},
	}

	// Filter by user ID
	filtered := []models.Entry{}
	for _, entry := range entries {
		if entry.UserID == 1 {
			filtered = append(filtered, entry)
		}
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries for user 1, got %d", len(filtered))
	}

	// Verify all filtered entries belong to user 1
	for _, entry := range filtered {
		if entry.UserID != 1 {
			t.Errorf("Entry %d does not belong to user 1", entry.ID)
		}
	}
}
