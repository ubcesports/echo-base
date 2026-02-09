package internal

import (
	"net/http"

	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/services"
)

func AddRoutes(
	mux *http.ServeMux,
	authService services.AuthService,
	gamerProfileService services.GamerProfileService,
	gamerActivityService services.GamerActivityService,
) {
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/db/ping", handlers.DatabasePing)
	mux.Handle("POST /admin/generate-key", handlers.GenerateAPIKey(authService))

	mux.Handle("GET /v1/api/gamer/{student_number}", handlers.GetGamerProfile(gamerProfileService))
	mux.Handle("POST /v1/api/gamer", handlers.CreateOrUpdateGamerProfile(gamerProfileService))
	mux.Handle("DELETE /v1/api/gamer/{student_number}", handlers.DeleteGamerProfile(gamerProfileService))

	mux.Handle("GET /v1/api/activity/{student_number}", handlers.GetActivityByStudent(gamerActivityService))
	mux.Handle("GET /v1/api/activity/today/{student_number}", handlers.GetTodayActivityByStudent(gamerActivityService))
	mux.Handle("GET /v1/api/activity/all/recent", handlers.GetRecentActivities(gamerActivityService))
	mux.Handle("POST /v1/api/activity", handlers.StartActivity(gamerActivityService))
	mux.Handle("PATCH /v1/api/activity/update/{student_number}", handlers.EndActivity(gamerActivityService))
	mux.Handle("GET /v1/api/activity/all/get-active-pcs", handlers.GetActiveSessions(gamerActivityService))
}
