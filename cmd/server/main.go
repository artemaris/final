package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/artemaris/gophkeeper/internal/server"
	"github.com/artemaris/gophkeeper/internal/storage"
)

var (
	addr     = flag.String("addr", ":8080", "Server address")
	tlsCert  = flag.String("tls-cert", "", "Path to TLS certificate file (enables HTTPS)")
	tlsKey   = flag.String("tls-key", "", "Path to TLS private key file (enables HTTPS)")
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

	// Get TLS configuration from flags or environment
	certFile := *tlsCert
	keyFile := *tlsKey
	
	// If not provided via flags, try environment variables
	if certFile == "" {
		certFile = os.Getenv("TLS_CERT")
	}
	if keyFile == "" {
		keyFile = os.Getenv("TLS_KEY")
	}

	// Start server with TLS if certificates are provided
	if certFile != "" && keyFile != "" {
		log.Printf("Starting HTTPS server on %s", *addr)
		log.Printf("TLS Certificate: %s", certFile)
		log.Printf("TLS Key: %s", keyFile)
		if err := http.ListenAndServeTLS(*addr, certFile, keyFile, srv); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server on %s (TLS not configured)", *addr)
		log.Printf("To enable HTTPS, provide --tls-cert and --tls-key flags or TLS_CERT and TLS_KEY environment variables")
		if err := http.ListenAndServe(*addr, srv); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}

func loadEnv() error {
	// Try to load .env file if it exists
	// If file doesn't exist, it's not an error - use environment variables or defaults
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return nil
	}
	// Load environment variables from .env file
	return godotenv.Load()
}
