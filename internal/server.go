package internal

import (
	"net/http"

	"github.com/ubcesports/echo-base/internal/middleware"
	"github.com/ubcesports/echo-base/internal/services"
)

func NewServer(
	authService services.AuthService,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		authService,
	)

	var handler http.Handler = mux
	handler = middleware.AuthMiddleware(handler, authService)

	return handler

}
