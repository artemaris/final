package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/artemaris/gophkeeper/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	db     storage.Storage
	router *http.ServeMux
}

// NewServer creates a new server instance
func NewServer(db storage.Storage) *Server {
	s := &Server{
		db:     db,
		router: http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Auth routes (public)
	s.router.HandleFunc("/api/register", s.RegisterHandler)
	s.router.HandleFunc("/api/login", s.LoginHandler)

	// Entry routes (protected)
	s.router.HandleFunc("/api/entries", s.AuthMiddleware(s.EntriesHandler))
	s.router.HandleFunc("/api/entries/", s.AuthMiddleware(s.EntryHandler))
	s.router.HandleFunc("/api/sync", s.AuthMiddleware(s.SyncHandler))
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Helper functions

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// getUserID retrieves the user ID from the request context
func getUserID(r *http.Request) (int, error) {
	ctx := r.Context()
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok || userID == 0 {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// getUserName retrieves the username from the request context
func getUserName(r *http.Request) (string, error) {
	ctx := r.Context()
	username, ok := ctx.Value(UsernameKey).(string)
	if !ok || username == "" {
		return "", fmt.Errorf("username not found in context")
	}
	return username, nil
}

func parseEntryID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[2])
}
