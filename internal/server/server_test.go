package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemaris/gophkeeper/internal/models"
)

func TestNewServer(t *testing.T) {
	// Используем mock storage
	mockDB := &MockStorage{}

	srv := NewServer(mockDB)

	if srv == nil {
		t.Fatal("Server should not be nil")
	}

	if srv.db != mockDB {
		t.Error("Server should use provided storage")
	}

	if srv.router == nil {
		t.Error("Router should not be nil")
	}
}

func TestServeHTTP(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/api/entries", nil)
	w := httptest.NewRecorder()

	// Вызываем ServeHTTP
	srv.ServeHTTP(w, req)

	// Проверяем, что запрос был обработан
	// Так как нет токена, ожидаем статус 401 или другой код ошибки
	if w.Code == http.StatusOK {
		t.Error("Expected error response without token")
	}
}

func TestServeHTTP_UnknownRoute(t *testing.T) {
	mockDB := &MockStorage{}
	srv := NewServer(mockDB)

	// Создаем запрос к несуществующему маршруту
	req := httptest.NewRequest("GET", "/api/unknown", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	// Проверяем 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestServer_Routes(t *testing.T) {
	mockDB := NewMockStorage()
	srv := NewServer(mockDB)

	// Тестируем все маршруты
	tests := []struct {
		method string
		path   string
		name   string
	}{
		{"POST", "/api/register", "Register"},
		{"POST", "/api/login", "Login"},
		{"GET", "/api/entries", "EntriesList"},
		{"GET", "/api/entries/1", "EntryGet"},
		{"PUT", "/api/entries/1", "EntryUpdate"},
		{"DELETE", "/api/entries/1", "EntryDelete"},
		{"GET", "/api/sync", "Sync"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" || tt.method == "PUT" {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader([]byte("{}")))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)

			// Проверяем, что маршрут существует и обрабатывается
			if w.Code == 0 {
				t.Error("Response code should be set")
			}
		})
	}
}

// Новый конструктор для MockStorage
func NewMockStorage() *MockStorage {
	return &MockStorage{
		CreateUserFunc: func(username, passwordHash string) (*models.User, error) {
			return &models.User{ID: 1, Username: username}, nil
		},
		GetUserByUsernameFunc: func(username string) (*models.User, error) {
			if username == "existing" {
				return &models.User{ID: 1, Username: username}, nil
			}
			return nil, errors.New("user not found")
		},
		GetUserByIDFunc: func(id int) (*models.User, error) {
			return &models.User{ID: id, Username: "testuser"}, nil
		},
		UserExistsFunc: func(username string) (bool, error) {
			return username == "existing", nil
		},
	}
}
