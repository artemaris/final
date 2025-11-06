package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8080")
	if c == nil {
		t.Fatal("Client should not be nil")
	}
	if c.baseURL != "http://localhost:8080" {
		t.Errorf("Expected baseURL 'http://localhost:8080', got '%s'", c.baseURL)
	}
}

func TestSetToken(t *testing.T) {
	c := NewClient("http://localhost:8080")
	token := "test-token"

	c.SetToken(token)

	if c.token != token {
		t.Errorf("Expected token '%s', got '%s'", token, c.token)
	}
}

func TestRegister_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/register" {
			t.Errorf("Expected path /api/register, got %s", r.URL.Path)
		}

		authResp := models.AuthResponse{
			Token: "test-token",
			User: models.User{
				ID:       1,
				Username: "testuser",
			},
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(authResp)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Register("testuser", "password123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected token in response")
	}

	if c.token != resp.Token {
		t.Error("Token should be set in client")
	}
}

func TestRegister_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Register("testuser", "password123")

	if err == nil {
		t.Error("Expected error")
	}
}

func TestLogin_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authResp := models.AuthResponse{
			Token: "login-token",
			User: models.User{
				ID:       1,
				Username: "testuser",
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authResp)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.Login("testuser", "password123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", resp.User.Username)
	}
}

func TestListEntries_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries := []models.Entry{
			{
				ID:    1,
				Type:  models.DataTypeLogin,
				Title: "Test Entry",
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entries)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	entries, err := c.ListEntries()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
}

func TestGetEntry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := models.Entry{
			ID:    1,
			Type:  models.DataTypeLogin,
			Title: "Test Entry",
			Data:  "test-data",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entry)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	entry, err := c.GetEntry(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.ID != 1 {
		t.Errorf("Expected ID 1, got %d", entry.ID)
	}
}

func TestCreateEntry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := models.Entry{
			ID:    1,
			Type:  models.DataTypeLogin,
			Title: "New Entry",
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entry)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.CreateEntryRequest{
		Type:  models.DataTypeLogin,
		Title: "New Entry",
		Data:  "test-data",
	}

	entry, err := c.CreateEntry(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.Title != "New Entry" {
		t.Errorf("Expected title 'New Entry', got '%s'", entry.Title)
	}
}

func TestDeleteEntry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	err := c.DeleteEntry(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSyncEntries_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		syncResp := models.SyncResponse{
			Entries: []models.Entry{
				{
					ID:   1,
					Type: models.DataTypeLogin,
				},
			},
			Version: int(time.Now().Unix()),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(syncResp)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	resp, err := c.SyncEntries(time.Unix(0, 0))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(resp.Entries))
	}
}

// Edge cases для client

func TestRegister_NetworkError(t *testing.T) {
	// Создаем клиент, указывающий на несуществующий сервер
	c := NewClient("http://localhost:99999")

	_, err := c.Register("testuser", "password123")
	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestRegister_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		// Empty body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Register("testuser", "password123")

	if err == nil {
		t.Error("Expected error for empty response")
	}
}

func TestLogin_NetworkError(t *testing.T) {
	c := NewClient("http://localhost:99999")

	_, err := c.Login("testuser", "password123")
	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestListEntries_EmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries := []models.Entry{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entries)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	entries, err := c.ListEntries()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}

func TestCreateEntry_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		// Empty body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.CreateEntryRequest{
		Type:  models.DataTypeLogin,
		Title: "New Entry",
		Data:  "test-data",
	}

	_, err := c.CreateEntry(req)

	if err == nil {
		t.Error("Expected error for empty response")
	}
}

func TestDeleteEntry_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "entry not found"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	err := c.DeleteEntry(1)

	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestDoRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := NewClient(server.URL)

	// Попытка отправить невалидный JSON
	_, err := c.doRequest("POST", "/test", make(chan int))

	if err == nil {
		t.Error("Expected error for invalid JSON body")
	}
}

// Дополнительные тесты для улучшения покрытия

func TestUpdateEntry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := models.Entry{
			ID:    1,
			Type:  models.DataTypeLogin,
			Title: "Updated Entry",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entry)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.UpdateEntryRequest{
		Title: "Updated Entry",
		Data:  "updated-data",
	}

	entry, err := c.UpdateEntry(1, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.Title != "Updated Entry" {
		t.Errorf("Expected title 'Updated Entry', got '%s'", entry.Title)
	}
}

func TestUpdateEntry_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "entry not found"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.UpdateEntryRequest{
		Title: "Updated Entry",
	}

	_, err := c.UpdateEntry(1, req)

	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestGetEntry_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "entry not found"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	_, err := c.GetEntry(999)

	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestListEntries_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("invalid-token")

	_, err := c.ListEntries()

	if err == nil {
		t.Error("Expected error for unauthorized")
	}
}

func TestCreateEntry_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.CreateEntryRequest{
		Type:  models.DataTypeLogin,
		Title: "",
	}

	_, err := c.CreateEntry(req)

	if err == nil {
		t.Error("Expected error for bad request")
	}
}

func TestSyncEntries_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	_, err := c.SyncEntries(time.Unix(0, 0))

	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestDoRequest_NetworkError(t *testing.T) {
	c := NewClient("http://localhost:99999")

	_, err := c.doRequest("GET", "/test", nil)

	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestGetEntry_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Empty body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	_, err := c.GetEntry(1)

	if err == nil {
		t.Error("Expected error for empty response")
	}
}

func TestUpdateEntry_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Empty body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.UpdateEntryRequest{
		Title: "Updated",
	}

	_, err := c.UpdateEntry(1, req)

	if err == nil {
		t.Error("Expected error for empty response")
	}
}

func TestSyncEntries_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Empty body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	_, err := c.SyncEntries(time.Unix(0, 0))

	if err == nil {
		t.Error("Expected error for empty response")
	}
}

func TestCreateEntry_StatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := models.Entry{
			ID:    1,
			Type:  models.DataTypeLogin,
			Title: "New Entry",
		}
		w.WriteHeader(http.StatusOK) // Some servers return 200 instead of 201
		json.NewEncoder(w).Encode(entry)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	c.SetToken("test-token")

	req := models.CreateEntryRequest{
		Type:  models.DataTypeLogin,
		Title: "New Entry",
		Data:  "test-data",
	}

	entry, err := c.CreateEntry(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.Title != "New Entry" {
		t.Errorf("Expected title 'New Entry', got '%s'", entry.Title)
	}
}

func TestRegister_EmptyErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		// No body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Register("testuser", "password123")

	// Should still return an error even without error message
	if err == nil {
		t.Error("Expected error")
	}
}

func TestLogin_EmptyErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		// No body
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Login("testuser", "password123")

	// Should still return an error even without error message
	if err == nil {
		t.Error("Expected error")
	}
}
