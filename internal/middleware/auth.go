package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

// TODO: This is just a dummy DB call, the real DB logic will live here
func getAPIKeyRecord(hash string) (userID string, err error) {
	validHash := "PUT HASHED KEY HERE"
	if hash != validHash {
		return "", errors.New("API key not found")
	}
	return "user_123", nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		hash := sha256.Sum256([]byte(apiKey))
		hashHex := hex.EncodeToString(hash[:])
		userID, err := getAPIKeyRecord(hashHex)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
