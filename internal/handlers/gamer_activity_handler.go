package handlers

import (
	"encoding/json"
	goerrors "errors"
	"net/http"
	"strconv"

	"github.com/ubcesports/echo-base/internal/errors"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/services"
)

func GetActivityByStudent(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		studentNumber := r.PathValue("student_number")

		activities, err := service.GetActivitiesByStudent(r.Context(), studentNumber)
		if err != nil {
			var validationErr *errors.ValidationError

			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(activities) == 0 {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(activities)
	})
}

func GetTodayActivityByStudent(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		studentNumber := r.PathValue("student_number")

		activities, err := service.GetTodayActivities(r.Context(), studentNumber)
		if err != nil {
			var validationErr *errors.ValidationError

			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(activities)
	})
}

func GetRecentActivities(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		page := 1
		limit := 10
		search := ""

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			var err error
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				http.Error(w, "Invalid page parameter", http.StatusBadRequest)
				return
			}
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
				return
			}
		}

		if searchStr := r.URL.Query().Get("search"); searchStr != "" {
			search = searchStr
		}

		activities, err := service.GetRecentActivities(r.Context(), page, limit, search)
		if err != nil {
			var validationErr *errors.ValidationError

			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(activities)
	})
}

func StartActivity(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.CreateActivityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		activity, err := service.StartActivity(r.Context(), &req)
		if err != nil {
			var validationErr *errors.ValidationError
			var notFoundErr *errors.NotFoundError
			var forbiddenErr *errors.ForbiddenError

			if goerrors.As(err, &notFoundErr) {
				http.Error(w, "Foreign key "+req.StudentNumber+" not found.", http.StatusNotFound)
				return
			}
			if goerrors.As(err, &forbiddenErr) {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(activity)
	})
}

func EndActivity(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		studentNumber := r.PathValue("student_number")

		var req models.UpdateActivityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		activity, err := service.EndActivity(r.Context(), studentNumber, &req)
		if err != nil {
			var validationErr *errors.ValidationError

			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if activity == nil {
				http.Error(w, "Student not active.", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(activity)
	})
}

func GetActiveSessions(service services.GamerActivityService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		activities, err := service.GetActiveSessions(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(activities)
	})
}
