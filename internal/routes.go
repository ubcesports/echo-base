package internal

import (
	"net/http"

	"github.com/ubcesports/echo-base/internal/handlers"
)

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/db/ping", handlers.DatabasePing)
}
