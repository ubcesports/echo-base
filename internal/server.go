package internal

import (
	"net/http"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	AddRoutes(mux)

	var handler http.Handler = mux
	//
	//	var authService auth.AuthService
	//	authMiddleware := middleware.NewAuthMiddleware(authService)
	//
	//	handler = authMiddleware.Middleware(handler)
	return handler

}
