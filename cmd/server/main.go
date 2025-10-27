package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/artemaris/gophkeeper/internal/server"
	"github.com/artemaris/gophkeeper/internal/storage"
)

var (
	addr = flag.String("addr", ":8080", "Server address")
)

func main() {
	flag.Parse()

	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	// Connect to database
	dsn := storage.GetDataSourceName()
	db, err := storage.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create server
	srv := server.NewServer(db)

	// Start server
	log.Printf("Starting server on %s", *addr)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loadEnv() error {
	// Try to load .env file if it exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return nil
	}
	return nil
}
