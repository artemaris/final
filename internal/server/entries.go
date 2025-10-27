package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// EntriesHandler handles listing and creating entries
func (s *Server) EntriesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.listEntries(w, r, userID)
	case http.MethodPost:
		s.createEntry(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// EntryHandler handles operations on a specific entry
func (s *Server) EntryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse entry ID from path
	entryID, err := parseEntryID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid entry ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getEntry(w, r, userID, entryID)
	case http.MethodPut:
		s.updateEntry(w, r, userID, entryID)
	case http.MethodDelete:
		s.deleteEntry(w, r, userID, entryID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listEntries(w http.ResponseWriter, r *http.Request, userID int) {
	entries, err := s.db.ListEntries(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list entries")
		return
	}

	respondJSON(w, http.StatusOK, entries)
}

func (s *Server) createEntry(w http.ResponseWriter, r *http.Request, userID int) {
	var req models.CreateEntryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate
	if req.Type == "" || req.Title == "" || req.Data == "" {
		respondError(w, http.StatusBadRequest, "Type, title, and data are required")
		return
	}

	entry, err := s.db.CreateEntry(userID, req.Type, req.Title, req.Data, req.Metadata)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create entry")
		return
	}

	respondJSON(w, http.StatusCreated, entry)
}

func (s *Server) getEntry(w http.ResponseWriter, r *http.Request, userID, entryID int) {
	entry, err := s.db.GetEntry(entryID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Entry not found")
		return
	}

	respondJSON(w, http.StatusOK, entry)
}

func (s *Server) updateEntry(w http.ResponseWriter, r *http.Request, userID, entryID int) {
	var req models.UpdateEntryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	entry, err := s.db.UpdateEntry(entryID, userID, req.Title, req.Data, req.Metadata, req.Version)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to update entry")
		return
	}

	respondJSON(w, http.StatusOK, entry)
}

func (s *Server) deleteEntry(w http.ResponseWriter, r *http.Request, userID, entryID int) {
	if err := s.db.DeleteEntry(entryID, userID); err != nil {
		respondError(w, http.StatusNotFound, "Entry not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Entry deleted successfully"})
}

// SyncHandler handles data synchronization
func (s *Server) SyncHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse since parameter
	sinceStr := r.URL.Query().Get("since")
	var since time.Time
	if sinceStr != "" {
		timestamp, err := strconv.ParseInt(sinceStr, 10, 64)
		if err == nil {
			since = time.Unix(timestamp, 0)
		}
	}

	entries, err := s.db.SyncEntries(userID, since)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to sync entries")
		return
	}

	respondJSON(w, http.StatusOK, models.SyncResponse{
		Entries: entries,
		Version: int(time.Now().Unix()),
	})
}
