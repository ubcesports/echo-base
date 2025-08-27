package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/ubcesports/echo-base/internal/database"
)

const (
	KeyIDLength      = 8
	SecretLength     = 32
	APIKeyPrefix     = "api_"
	MaxAppNameLength = 100
)

var (
	validAppNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type GenerateKeyRequest struct {
	AppName string `json:"app_name"`
}

type GenerateKeyResponse struct {
	KeyID   string `json:"key_id"`
	APIKey  string `json:"api_key"`
	AppName string `json:"app_name"`
}

func GenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateAppName(req.AppName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyIDBytes := make([]byte, KeyIDLength)
	_, err := rand.Read(keyIDBytes)
	if err != nil {
		http.Error(w, "Failed to generate key ID", http.StatusInternalServerError)
		return
	}
	keyID := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(keyIDBytes))

	secretBytes := make([]byte, SecretLength)
	_, err = rand.Read(secretBytes)
	if err != nil {
		http.Error(w, "Failed to generate secret", http.StatusInternalServerError)
		return
	}
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedSecret := hasher.Sum(nil)

	if err := storeAPIKey(req, keyID, hashedSecret); err != nil {
		http.Error(w, "Failed to store API Key", http.StatusInternalServerError)
		return
	}

	fullKey := fmt.Sprintf("api_%s.%s", keyID, secret)
	response := GenerateKeyResponse{
		KeyID:   keyID,
		APIKey:  fullKey,
		AppName: req.AppName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func storeAPIKey(req GenerateKeyRequest, keyID string, hashedSecret []byte) error {
	query := `
        INSERT INTO application (app_name, key_id, hashed_key)
        VALUES ($1, $2, $3)
    `
	_, err := database.DB.Exec(query, req.AppName, keyID, hashedSecret)
	if err != nil {
		return fmt.Errorf("database storage failed")

	}

	return nil
}

func validateAppName(appName string) error {
	if appName == "" {
		return fmt.Errorf("app_name is required")
	}

	if len(appName) > MaxAppNameLength {
		return fmt.Errorf("app_name must be %d characters or less", MaxAppNameLength)
	}

	if !validAppNameRegex.MatchString(appName) {
		return fmt.Errorf("app_name can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}
