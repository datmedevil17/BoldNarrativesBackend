package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/datmedevil17/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get port from environment variable with default fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	fmt.Printf("Server is starting on port %s...\n", port)

	r := router.Router()
	log.Fatal(http.ListenAndServe(":"+port, r))
}
