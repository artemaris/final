package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
	"github.com/artemaris/gophkeeper/internal/client"
	"github.com/artemaris/gophkeeper/internal/models"
)

var (
	version    = "dev"
	buildDate  = "unknown"
	serverURL  = flag.String("server", "", "Server URL (e.g., http://localhost:8080 or https://api.example.com)")
	timeoutSec = flag.Int("timeout", 0, "Request timeout in seconds (default: 30, or CLIENT_TIMEOUT env var)")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[0]
	switch cmd {
	case "version":
		fmt.Printf("GophKeeper Client\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Date: %s\n", buildDate)
	case "register":
		handleRegister()
	case "login":
		handleLogin()
	case "add":
		handleAdd()
	case "list":
		handleList()
	case "get":
		handleGet()
	case "sync":
		handleSync()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("GophKeeper - Password Manager CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gophkeeper-client [flags] <command> [arguments]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -server string   Server URL (default: http://localhost:8080)")
	fmt.Println("                   Can also be set via SERVER_URL environment variable")
	fmt.Println("  -timeout int     Request timeout in seconds (default: 30)")
	fmt.Println("                   Can also be set via CLIENT_TIMEOUT environment variable")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version          Show version information")
	fmt.Println("  register         Register a new user")
	fmt.Println("  login            Login to the server")
	fmt.Println("  add <type>       Add a new entry (login/text/binary/card)")
	fmt.Println("  list             List all entries")
	fmt.Println("  get <id>         Get an entry by ID")
	fmt.Println("  sync             Sync with server")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gophkeeper-client -server http://localhost:8080 register")
	fmt.Println("  gophkeeper-client -server https://api.example.com login")
	fmt.Println("  gophkeeper-client -timeout 120 add binary  # Upload large files with 2 minute timeout")
	fmt.Println("  SERVER_URL=https://api.example.com gophkeeper-client list")
	fmt.Println("  CLIENT_TIMEOUT=60 gophkeeper-client sync   # Use 60 second timeout")
}

func handleRegister() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	password, err := readPassword("Password: ")
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}

	c := newClient()
	resp, err := c.Register(username, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nRegistered successfully! User ID: %d\n", resp.User.ID)
	fmt.Printf("Token: %s\n", resp.Token)
}

func handleLogin() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	password, err := readPassword("Password: ")
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}

	c := newClient()
	resp, err := c.Login(username, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nLogged in successfully! User ID: %d\n", resp.User.ID)
	fmt.Printf("Token: %s\n", resp.Token)
}

func handleAdd() {
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: add <type>")
		fmt.Println("Types: login, text, binary, card")
		os.Exit(1)
	}

	entryType := args[1]
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Data: ")
	scanner.Scan()
	data := scanner.Text()

	fmt.Print("Metadata (optional): ")
	scanner.Scan()
	metadata := scanner.Text()

	c := newClient()
	// Note: In production, this would read the token from a secure store
	// For now, you need to set it manually after login
	fmt.Println("Warning: Token not set. Please login first and set token.")

	req := models.CreateEntryRequest{
		Type:     models.DataType(entryType),
		Title:    title,
		Data:     data,
		Metadata: metadata,
	}

	entry, err := c.CreateEntry(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Entry created successfully! ID: %d\n", entry.ID)
}

func handleList() {
	c := newClient()

	entries, err := c.ListEntries()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	fmt.Printf("Found %d entries:\n\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("ID: %d | Type: %s | Title: %s\n", entry.ID, entry.Type, entry.Title)
	}
}

func handleGet() {
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: get <id>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Invalid ID: %s\n", args[1])
		os.Exit(1)
	}

	c := newClient()
	entry, err := c.GetEntry(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID: %d\n", entry.ID)
	fmt.Printf("Type: %s\n", entry.Type)
	fmt.Printf("Title: %s\n", entry.Title)
	fmt.Printf("Data: %s\n", entry.Data)
	fmt.Printf("Metadata: %s\n", entry.Metadata)
	fmt.Printf("Version: %d\n", entry.Version)
	fmt.Printf("Created: %s\n", entry.CreatedAt)
}

func handleSync() {
	c := newClient()

	// Sync entries from the beginning of time
	since := time.Unix(0, 0)
	resp, err := c.SyncEntries(since)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Synced %d entries. Version: %d\n", len(resp.Entries), resp.Version)
}

// getServerURL returns the server URL from command line flag, environment variable, or default value
func getServerURL() string {
	// First, try command line flag
	if *serverURL != "" {
		return *serverURL
	}
	
	// Then, try environment variable
	if url := os.Getenv("SERVER_URL"); url != "" {
		return url
	}
	
	// Default value
	return "http://localhost:8080"
}

// getTimeout returns the timeout from command line flag, environment variable, or default value
func getTimeout() time.Duration {
	// First, try command line flag
	if *timeoutSec > 0 {
		return time.Duration(*timeoutSec) * time.Second
	}
	
	// Then, try environment variable (handled by client package)
	// If not set, client package will use default (30 seconds)
	return 0
}

// newClient creates a new client with configured server URL and timeout
func newClient() *client.Client {
	timeout := getTimeout()
	if timeout > 0 {
		return client.NewClientWithTimeout(getServerURL(), timeout)
	}
	return client.NewClient(getServerURL())
}

// readPassword securely reads a password from stdin without echoing it to the terminal
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	
	// Read password byte-by-byte from stdin without echoing
	// Use os.Stdin.Fd() for cross-platform compatibility
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	
	// Print newline after password input (since it wasn't printed automatically)
	fmt.Println()
	
	return strings.TrimSpace(string(passwordBytes)), nil
}
