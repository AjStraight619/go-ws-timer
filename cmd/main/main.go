package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AjStraight619/go-ws-timer/internal/routes"
)

func main() {
	r := routes.SetUpRouter()
	routes.SetUpRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no PORT environment variable is set
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
