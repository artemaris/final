package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"username":"test","password":"pass"}`))
		req.Header.Set("Content-Type", "application/json")

		var result struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := decodeJSON(req, &result)
		if err != nil {
			t.Fatalf("Failed to decode JSON: %v", err)
		}

		if result.Username != "test" {
			t.Errorf("Expected username 'test', got '%s'", result.Username)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		var result map[string]interface{}

		err := decodeJSON(req, &result)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	respondJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got '%s'", result["message"])
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()

	respondError(w, http.StatusBadRequest, "Test error")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "Test error" {
		t.Errorf("Expected error 'Test error', got '%s'", result["error"])
	}
}

func TestGetUserID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(r.Context(), UserIDKey, 42)
	r = r.WithContext(ctx)

	userID, err := getUserID(r)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	if userID != 42 {
		t.Errorf("Expected user ID 42, got %d", userID)
	}
}

func TestGetUserID_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	_, err := getUserID(r)
	if err == nil {
		t.Error("Expected error when user ID is missing from context")
	}
}

func TestGetUserName(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(r.Context(), UsernameKey, "testuser")
	r = r.WithContext(ctx)

	username, err := getUserName(r)
	if err != nil {
		t.Fatalf("Failed to get username: %v", err)
	}

	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
}

func TestGetUserName_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	_, err := getUserName(r)
	if err == nil {
		t.Error("Expected error when username is missing from context")
	}
}

func TestParseEntryID(t *testing.T) {
	tests := []struct {
		path     string
		expected int
		err      bool
	}{
		{"/api/entries/1", 1, false},
		{"/api/entries/42", 42, false},
		{"/api/entries/invalid", 0, true},
		{"/api/entries", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			id, err := parseEntryID(tt.path)
			if tt.err && err == nil {
				t.Error("Expected error")
			}
			if !tt.err && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !tt.err && id != tt.expected {
				t.Errorf("Expected ID %d, got %d", tt.expected, id)
			}
		})
	}
}
