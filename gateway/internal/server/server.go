package server

import (
	"fmt"
	"net/http"

	"gateway/internal/handlers"
)

// Start initializes and starts the HTTP server
func Start() error {
	// Register HTTP handlers
	http.HandleFunc("/readPP/", handlers.ReadPPHandler)

	// Define the port number
	port := ":8080"

	// Print a message indicating the server is starting
	fmt.Printf("\nServer listening on port %s...\n", port)

	// Start the HTTP server
	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}
	return nil
}
