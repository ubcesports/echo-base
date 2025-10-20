package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

type AuthMiddleware struct {
	authService auth.AuthService
}

func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		appName, err := m.authService.ValidateAPIKey(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "appName", appName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
