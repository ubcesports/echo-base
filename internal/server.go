package internal

import (
	"net/http"

	"github.com/ubcesports/echo-base/internal/middleware"
	"github.com/ubcesports/echo-base/internal/services"
)

func NewServer(
	authService services.AuthService,
	gamerProfileService services.GamerProfileService,
	gamerActivityService services.GamerActivityService,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		authService,
		gamerProfileService,
		gamerActivityService,
	)

	var handler http.Handler = mux
	handler = middleware.AuthMiddleware(handler, authService)

	return handler

}
