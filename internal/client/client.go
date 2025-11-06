package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// Client represents the GophKeeper client
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new GophKeeper client
// If timeout is 0, it will be read from CLIENT_TIMEOUT environment variable
// or default to 30 seconds if the environment variable is not set
func NewClient(baseURL string) *Client {
	return NewClientWithTimeout(baseURL, 0)
}

// NewClientWithTimeout creates a new GophKeeper client with a custom timeout
// If timeout is 0, it will be read from CLIENT_TIMEOUT environment variable
// or default to 30 seconds if the environment variable is not set
func NewClientWithTimeout(baseURL string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = getTimeoutFromEnv()
	}

	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// getTimeoutFromEnv reads the timeout from CLIENT_TIMEOUT environment variable
// Returns default timeout (30 seconds) if the variable is not set or invalid
func getTimeoutFromEnv() time.Duration {
	timeoutStr := os.Getenv("CLIENT_TIMEOUT")
	if timeoutStr == "" {
		return 30 * time.Second
	}

	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSeconds <= 0 {
		// Invalid value, return default
		return 30 * time.Second
	}

	return time.Duration(timeoutSeconds) * time.Second
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// Register registers a new user
func (c *Client) Register(username, password string) (*models.AuthResponse, error) {
	req := models.RegisterRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.doRequest(http.MethodPost, "/api/register", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("registration failed: %v", errResp["error"])
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	c.token = authResp.Token
	return &authResp, nil
}

// Login authenticates a user
func (c *Client) Login(username, password string) (*models.AuthResponse, error) {
	req := models.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.doRequest(http.MethodPost, "/api/login", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("login failed: %v", errResp["error"])
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	c.token = authResp.Token
	return &authResp, nil
}

// ListEntries retrieves all entries
func (c *Client) ListEntries() ([]models.Entry, error) {
	resp, err := c.doRequest(http.MethodGet, "/api/entries", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("failed to list entries: %v", errResp["error"])
	}

	var entries []models.Entry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetEntry retrieves an entry by ID
func (c *Client) GetEntry(id int) (*models.Entry, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/entries/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("failed to get entry: %v", errResp["error"])
	}

	var entry models.Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// CreateEntry creates a new entry
func (c *Client) CreateEntry(entry models.CreateEntryRequest) (*models.Entry, error) {
	resp, err := c.doRequest(http.MethodPost, "/api/entries", entry)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("failed to create entry: %v", errResp["error"])
	}

	var createdEntry models.Entry
	if err := json.NewDecoder(resp.Body).Decode(&createdEntry); err != nil {
		return nil, err
	}

	return &createdEntry, nil
}

// UpdateEntry updates an existing entry
func (c *Client) UpdateEntry(id int, entry models.UpdateEntryRequest) (*models.Entry, error) {
	resp, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/entries/%d", id), entry)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("failed to update entry: %v", errResp["error"])
	}

	var updatedEntry models.Entry
	if err := json.NewDecoder(resp.Body).Decode(&updatedEntry); err != nil {
		return nil, err
	}

	return &updatedEntry, nil
}

// DeleteEntry deletes an entry
func (c *Client) DeleteEntry(id int) error {
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/entries/%d", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("failed to delete entry: %v", errResp["error"])
	}

	return nil
}

// SyncEntries synchronizes entries with the server
func (c *Client) SyncEntries(since time.Time) (*models.SyncResponse, error) {
	url := fmt.Sprintf("/api/sync?since=%d", since.Unix())
	resp, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("failed to sync entries: %v", errResp["error"])
	}

	var syncResp models.SyncResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncResp); err != nil {
		return nil, err
	}

	return &syncResp, nil
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.client.Do(req)
}
