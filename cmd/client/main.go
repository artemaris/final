package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/artemaris/gophkeeper/internal/client"
	"github.com/artemaris/gophkeeper/internal/models"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

const (
	serverURL = "http://localhost:8080"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
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
	fmt.Println("  gophkeeper-client <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version          Show version information")
	fmt.Println("  register         Register a new user")
	fmt.Println("  login            Login to the server")
	fmt.Println("  add <type>       Add a new entry (login/text/binary/card)")
	fmt.Println("  list             List all entries")
	fmt.Println("  get <id>         Get an entry by ID")
	fmt.Println("  sync             Sync with server")
}

func handleRegister() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	c := client.NewClient(serverURL)
	resp, err := c.Register(username, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Registered successfully! User ID: %d\n", resp.User.ID)
	fmt.Printf("Token: %s\n", resp.Token)
}

func handleLogin() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	c := client.NewClient(serverURL)
	resp, err := c.Login(username, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Logged in successfully! User ID: %d\n", resp.User.ID)
	fmt.Printf("Token: %s\n", resp.Token)
}

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: add <type>")
		fmt.Println("Types: login, text, binary, card")
		os.Exit(1)
	}

	entryType := os.Args[2]
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

	c := client.NewClient(serverURL)
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
	c := client.NewClient(serverURL)

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
	if len(os.Args) < 3 {
		fmt.Println("Usage: get <id>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid ID: %s\n", os.Args[2])
		os.Exit(1)
	}

	c := client.NewClient(serverURL)
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
	c := client.NewClient(serverURL)

	// Sync entries from the beginning of time
	since := time.Unix(0, 0)
	resp, err := c.SyncEntries(since)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Synced %d entries. Version: %d\n", len(resp.Entries), resp.Version)
}
