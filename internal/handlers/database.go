package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ubcesports/echo-base/internal/database"
)

type DatabasePingResponse struct {
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
	ResponseTime string    `json:"response_time"`
	Error        string    `json:"error,omitempty"`
}

func DatabasePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := time.Now()

	response := DatabasePingResponse{
		Timestamp: time.Now(),
	}

	// Ping the database
	if err := database.Ping(); err != nil {
		response.Status = "error"
		response.Error = err.Error()
		response.ResponseTime = time.Since(start).String()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Status = "ok"
		response.ResponseTime = time.Since(start).String()
	}

	json.NewEncoder(w).Encode(response)
}
