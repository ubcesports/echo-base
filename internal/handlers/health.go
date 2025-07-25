package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ubcesports/echo-base/internal/database"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Database:  "ok",
	}

	// Check database connection
	if err := database.Ping(); err != nil {
		response.Status = "error"
		response.Database = "error: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}
