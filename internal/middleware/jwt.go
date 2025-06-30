package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	secret := os.Getenv("JWT_TOKEN")
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(writer, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}