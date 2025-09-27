package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/middleware"
)

func main() {
	config.LoadEnv(".env")
	cfg := config.LoadConfig()

	database.Init()
	defer database.Close()

	// Gamer Activity Handler
	ah := &handlers.Handler{
		Config: cfg,
		DB:     database.DB,
	}

	// Logger
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	slog.SetDefault(logger)

	// Initialize database connection
	database.Init()
	defer database.Close()

	// Set up HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/db/ping", handlers.DatabasePing)
	mux.HandleFunc("/admin/generate-key", handlers.GenerateAPIKey)
	mux.HandleFunc("/activity/{student_number}", handlers.Wrap(ah.GetGamerActivityByStudent))
	mux.HandleFunc("/activity", handlers.Wrap(ah.AddGamerActivity))
	mux.HandleFunc("/activity/today/{student_number}", handlers.Wrap(ah.GetGamerActivityByTierOneStudentToday))
	mux.HandleFunc("/activity/all/recent", handlers.Wrap(ah.GetGamerActivity))
	mux.HandleFunc("/activity/update/{student_number}", handlers.Wrap(ah.UpdateGamerActivity))
	mux.HandleFunc("/activity/get-active-pcs", handlers.Wrap(ah.GetAllActivePCs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Echo Base API is running!"))
	})

	// Apply auth middleware to all routes
	handler := middleware.AuthMiddleware(mux)

	// Get port from environment or default to 8080
	port := os.Getenv("EB_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at http://localhost:%s/health", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
