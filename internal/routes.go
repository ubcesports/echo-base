package internal

import (
	"net/http"

	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/services"
)

func AddRoutes(
	mux *http.ServeMux,
	authService services.AuthService,
) {
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/db/ping", handlers.DatabasePing)
	mux.Handle("/admin/generate-key", handlers.GenerateAPIKey(authService))
}
