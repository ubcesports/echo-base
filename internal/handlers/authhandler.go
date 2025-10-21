package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ubcesports/echo-base/internal/services"
)

type GenerateKeyRequest struct {
	AppName string `json:"app_name"`
}

func GenerateAPIKey(authService *services.AuthService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req GenerateKeyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			apiKey, err := authService.GenerateAPIKey(r.Context(), req.AppName)
			if err != nil {
				http.Error(w, "Error generating API key", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(apiKey)
		},
	)
}
