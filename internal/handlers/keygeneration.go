package handlers

import (
	"context"
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
	KeyIDLength      = 6
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

func CreateApiKey(appName string) (res GenerateKeyResponse, err error) {
	keyIDBytes := make([]byte, KeyIDLength)
	_, err = rand.Read(keyIDBytes)
	if err != nil {
		return GenerateKeyResponse{}, err
	}
	keyID := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(keyIDBytes))

	secretBytes := make([]byte, SecretLength)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return GenerateKeyResponse{}, err
	}
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedSecret := hasher.Sum(nil)

	if err := storeAPIKey(appName, keyID, hashedSecret); err != nil {
		return GenerateKeyResponse{}, err
	}

	fullKey := fmt.Sprintf("api_%s.%s", keyID, secret)

	response := GenerateKeyResponse{
		KeyID:   keyID,
		APIKey:  fullKey,
		AppName: appName,
	}

	return response, nil
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

	response, err := CreateApiKey(req.AppName)
	if err != nil {
		http.Error(w, "Error while generating API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func storeAPIKey(appName string, keyID string, hashedSecret []byte) error {
	query := `
        INSERT INTO application (app_name, key_id, hashed_key)
        	VALUES ($1, $2, $3)
    `
	_, err := database.DB.Exec(context.Background(), query, appName, keyID, hashedSecret)
	if err != nil {
		return fmt.Errorf("database storage failed: %s", err.Error())
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
