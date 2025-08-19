package handlers

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubcesports/echo-base/internal/database"
	"golang.org/x/crypto/argon2"
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
	var req GenerateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keyIDBytes := make([]byte, 6)
	_, err := rand.Read(keyIDBytes)
	if err != nil {
		return
	}
	keyID := strings.ToLower(base64.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(keyIDBytes))

	secretBytes := make([]byte, 32)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return
	}
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}
	hash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)
	hashedSecret := base64.RawURLEncoding.EncodeToString(salt) + ":" + base64.RawURLEncoding.EncodeToString(hash)

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

func storeAPIKey(req GenerateKeyRequest, keyID string, hashedSecret string) error {
	query := `
		INSERT INTO auth (app_name, key_id, hashed_key)
		VALUES ($1, $2, $3)
	`

	_, err := database.DB.Exec(query, req.AppName, keyID, hashedSecret)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
