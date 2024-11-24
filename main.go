package main

import (
	"log"
	"net/http"
	"os"

	"gitglyph/functions/webhook"
)

func main() {
	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}

	// Setup webhook route
	http.HandleFunc("/webhook", webhook.Handle)

	// Start the server
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
