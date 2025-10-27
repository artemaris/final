package storage

import (
	"testing"

	_ "github.com/lib/pq"
)

// TestGetDataSourceName tests database connection string generation
func TestGetDataSourceName(t *testing.T) {
	dsn := GetDataSourceName()
	if dsn == "" {
		t.Fatal("DataSourceName is empty")
	}

	// Should contain required fields
	if !contains(dsn, "host=") {
		t.Error("DataSourceName should contain host")
	}
	if !contains(dsn, "port=") {
		t.Error("DataSourceName should contain port")
	}
	if !contains(dsn, "user=") {
		t.Error("DataSourceName should contain user")
	}
	if !contains(dsn, "password=") {
		t.Error("DataSourceName should contain password")
	}
	if !contains(dsn, "dbname=") {
		t.Error("DataSourceName should contain dbname")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
