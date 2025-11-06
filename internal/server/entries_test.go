package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// requestWithUserContext creates a request with user context set for testing
func requestWithUserContext(r *http.Request, userID int, username string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, UserIDKey, userID)
	if username != "" {
		ctx = context.WithValue(ctx, UsernameKey, username)
	}
	return r.WithContext(ctx)
}

// Test error
var ErrUserNotFound = errors.New("user not found")

func TestCreateEntry(t *testing.T) {
	mockDB := &MockStorage{
		CreateEntryFunc: func(userID int, entryType models.DataType, title, data, metadata string) (*models.Entry, error) {
			return &models.Entry{
				ID:        1,
				UserID:    userID,
				Type:      entryType,
				Title:     title,
				Data:      data,
				Metadata:  metadata,
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	req := models.CreateEntryRequest{
		Type:     models.DataTypeLogin,
		Title:    "Test Entry",
		Data:     "encrypted_data",
		Metadata: "{\"site\":\"example.com\"}",
	}

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/api/entries", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.createEntry(response, request, 1)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, response.Code)
	}
}

func TestListEntries(t *testing.T) {
	mockDB := &MockStorage{
		ListEntriesFunc: func(userID int) ([]models.Entry, error) {
			return []models.Entry{
				{ID: 1, UserID: userID, Type: models.DataTypeLogin, Title: "Entry 1"},
				{ID: 2, UserID: userID, Type: models.DataTypeText, Title: "Entry 2"},
			}, nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	request := httptest.NewRequest(http.MethodGet, "/api/entries", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.listEntries(response, request, 1)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	var entries []models.Entry
	json.NewDecoder(response.Body).Decode(&entries)

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestGetEntry(t *testing.T) {
	mockDB := &MockStorage{
		GetEntryFunc: func(id, userID int) (*models.Entry, error) {
			return &models.Entry{
				ID:        1,
				UserID:    userID,
				Type:      models.DataTypeLogin,
				Title:     "Test Entry",
				Data:      "encrypted_data",
				Metadata:  "{}",
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	request := httptest.NewRequest(http.MethodGet, "/api/entries/1", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.getEntry(response, request, 1, 1)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestUpdateEntry(t *testing.T) {
	mockDB := &MockStorage{
		UpdateEntryFunc: func(id, userID int, title, data, metadata string, version int) (*models.Entry, error) {
			return &models.Entry{
				ID:        id,
				UserID:    userID,
				Type:      models.DataTypeLogin,
				Title:     title,
				Data:      data,
				Metadata:  metadata,
				Version:   version + 1,
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	req := models.UpdateEntryRequest{
		Title:    "Updated Entry",
		Data:     "updated_data",
		Metadata: "{\"site\":\"updated.com\"}",
		Version:  1,
	}

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPut, "/api/entries/1", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.updateEntry(response, request, 1, 1)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestDeleteEntry(t *testing.T) {
	mockDB := &MockStorage{
		DeleteEntryFunc: func(id, userID int) error {
			return nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	request := httptest.NewRequest(http.MethodDelete, "/api/entries/1", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.deleteEntry(response, request, 1, 1)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestSyncEntries(t *testing.T) {
	mockDB := &MockStorage{
		SyncEntriesFunc: func(userID int, since time.Time) ([]models.Entry, error) {
			return []models.Entry{
				{ID: 1, UserID: userID, Type: models.DataTypeLogin, Title: "Synced Entry"},
			}, nil
		},
	}

	srv := &Server{db: mockDB, router: http.NewServeMux()}

	request := httptest.NewRequest(http.MethodGet, "/api/sync?since=1640995200", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.SyncHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

// Edge cases для entries handlers

func TestCreateEntry_InvalidJSON(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	body := bytes.NewReader([]byte(`{"type":"login",`)) // Invalid JSON
	req := httptest.NewRequest("POST", "/api/entries", body)
	req.Header.Set("Content-Type", "application/json")
	req = requestWithUserContext(req, 1, "testuser")
	w := httptest.NewRecorder()

	srv.EntriesHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestCreateEntry_MissingFields(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	req := models.CreateEntryRequest{
		Type:     models.DataTypeLogin,
		Title:    "", // Empty title
		Data:     "data",
		Metadata: "",
	}

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/api/entries", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.createEntry(response, request, 1)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", response.Code)
	}
}

func TestCreateEntry_DatabaseError(t *testing.T) {
	mockDB := &MockStorage{
		CreateEntryFunc: func(userID int, entryType models.DataType, title, data, metadata string) (*models.Entry, error) {
			return nil, errors.New("database error")
		},
	}
	srv := NewServer(mockDB)

	req := models.CreateEntryRequest{
		Type:     models.DataTypeLogin,
		Title:    "Test",
		Data:     "data",
		Metadata: "",
	}

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/api/entries", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.createEntry(response, request, 1)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", response.Code)
	}
}

func TestUpdateEntry_InvalidJSON(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	body := bytes.NewReader([]byte(`{"title":"test",`)) // Invalid JSON
	req := httptest.NewRequest("PUT", "/api/entries/1", body)
	req.Header.Set("Content-Type", "application/json")
	req = requestWithUserContext(req, 1, "testuser")
	w := httptest.NewRecorder()

	srv.EntryHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestGetEntry_NotFound(t *testing.T) {
	mockDB := &MockStorage{
		GetEntryFunc: func(id, userID int) (*models.Entry, error) {
			return nil, errors.New("entry not found")
		},
	}
	srv := NewServer(mockDB)

	request := httptest.NewRequest(http.MethodGet, "/api/entries/999", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.getEntry(response, request, 1, 999)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", response.Code)
	}
}

func TestDeleteEntry_NotFound(t *testing.T) {
	mockDB := &MockStorage{
		DeleteEntryFunc: func(id, userID int) error {
			return errors.New("entry not found")
		},
	}
	srv := NewServer(mockDB)

	request := httptest.NewRequest(http.MethodDelete, "/api/entries/999", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.deleteEntry(response, request, 1, 999)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", response.Code)
	}
}

func TestListEntries_DatabaseError(t *testing.T) {
	mockDB := &MockStorage{
		ListEntriesFunc: func(userID int) ([]models.Entry, error) {
			return nil, errors.New("database error")
		},
	}
	srv := NewServer(mockDB)

	request := httptest.NewRequest(http.MethodGet, "/api/entries", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.listEntries(response, request, 1)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", response.Code)
	}
}

func TestSyncEntries_DatabaseError(t *testing.T) {
	mockDB := &MockStorage{
		SyncEntriesFunc: func(userID int, since time.Time) ([]models.Entry, error) {
			return nil, errors.New("database error")
		},
	}
	srv := NewServer(mockDB)

	request := httptest.NewRequest(http.MethodGet, "/api/sync?since=1640995200", nil)
	request = requestWithUserContext(request, 1, "testuser")
	response := httptest.NewRecorder()

	srv.SyncHandler(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", response.Code)
	}
}

func TestEntryHandler_InvalidEntryID(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	req := httptest.NewRequest("GET", "/api/entries/invalid", nil)
	req = requestWithUserContext(req, 1, "testuser")
	w := httptest.NewRecorder()

	srv.EntryHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestEntriesHandler_MethodNotAllowed(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	req := httptest.NewRequest("PUT", "/api/entries", nil)
	req = requestWithUserContext(req, 1, "testuser")
	w := httptest.NewRecorder()

	srv.EntriesHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}
}
