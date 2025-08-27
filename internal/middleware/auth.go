package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/ubcesports/echo-base/internal/database"
)

func verifySHA256(secret string, hashedSecret []byte) bool {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	actualHash := hasher.Sum(nil)

	return subtle.ConstantTimeCompare(hashedSecret, actualHash) == 1
}

func processAPIKey(apiKey string) (appName string, err error) {
	if !strings.HasPrefix(apiKey, "api_") {
		return "", sql.ErrNoRows
	}

	keyParts := strings.Split(strings.TrimPrefix(apiKey, "api_"), ".")
	if len(keyParts) != 2 {
		return "", sql.ErrNoRows
	}

	keyID := keyParts[0]
	secret := keyParts[1]

	var hashedSecret []byte
	var lastUsed *time.Time

	query := `
        SELECT app_name, hashed_key, last_used_at 
        FROM application 
        WHERE key_id = $1
    `
	err = database.DB.QueryRow(query, keyID).Scan(&appName, &hashedSecret, &lastUsed)
	if err != nil {
		return "", err
	}

	if !verifySHA256(secret, hashedSecret) {
		return "", sql.ErrNoRows
	}

	updateQuery := `UPDATE application SET last_used_at = NOW() WHERE key_id = $1`
	go func() {
		database.DB.Exec(updateQuery, keyID)
	}()

	return appName, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		appName, err := processAPIKey(apiKey)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "appName", appName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
