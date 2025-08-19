package middleware

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ubcesports/echo-base/internal/database"
	"golang.org/x/crypto/argon2"
)

func verifyArgon2(secret, hashedSecret string) bool {
	parts := strings.Split(hashedSecret, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	actualHash := argon2.IDKey([]byte(secret), salt, 1, 64*1024, 4, 32)

	if len(expectedHash) != len(actualHash) {
		return false
	}
	for i := range expectedHash {
		if expectedHash[i] != actualHash[i] {
			return false
		}
	}
	return true
}

func getAPIKeyRecord(apiKey string) (appName string, err error) {
	if !strings.HasPrefix(apiKey, "api_") {
		return "", sql.ErrNoRows
	}

	keyParts := strings.Split(strings.TrimPrefix(apiKey, "api_"), ".")
	if len(keyParts) != 2 {
		return "", sql.ErrNoRows
	}

	keyID := keyParts[0]
	secret := keyParts[1]

	var hashedSecret string
	var lastUsed *time.Time

	query := `
        SELECT app_name, hashed_key, last_used_at 
        FROM auth 
        WHERE key_id = $1
    `

	err = database.DB.QueryRow(query, keyID).Scan(&appName, &hashedSecret, &lastUsed)
	if err != nil {
		return "", err
	}

	if !verifyArgon2(secret, hashedSecret) {
		return "", sql.ErrNoRows
	}

	updateQuery := `UPDATE auth SET last_used_at = NOW() WHERE key_id = $1`
	database.DB.Exec(updateQuery, keyID)

	return appName, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("here")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		appName, err := getAPIKeyRecord(apiKey)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "appName", appName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
