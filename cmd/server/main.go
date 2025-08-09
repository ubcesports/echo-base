package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/middleware"
)

func main() {
	config.LoadEnv(".env")

	// Initialize database connection
	database.Init()
	defer database.Close()

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/health", middleware.AuthMiddleware(http.HandlerFunc(handlers.HealthCheck)))
	mux.Handle("/db/ping", middleware.AuthMiddleware(http.HandlerFunc(handlers.DatabasePing)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Echo Base API is running!"))
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at http://localhost:%s/health", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
