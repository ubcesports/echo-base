package handlers

import (
	"encoding/json"
	goerrors "errors"
	"net/http"

	"github.com/ubcesports/echo-base/internal/errors"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/services"
)

func GetGamerProfile(service services.GamerProfileService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		studentNumber := r.PathValue("student_number")

		profile, err := service.GetProfile(r.Context(), studentNumber)
		if err != nil {
			var notFoundErr *errors.NotFoundError
			var validationErr *errors.ValidationError

			if goerrors.As(err, &notFoundErr) {
				http.Error(w, "Student not found", http.StatusNotFound)
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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(profile)
	})
}

func CreateOrUpdateGamerProfile(service services.GamerProfileService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.CreateGamerProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		profile, err := service.CreateOrUpdateProfile(r.Context(), &req)
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
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(profile)
	})
}

func DeleteGamerProfile(service services.GamerProfileService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		studentNumber := r.PathValue("student_number")

		err := service.DeleteProfile(r.Context(), studentNumber)
		if err != nil {
			var notFoundErr *errors.NotFoundError
			var validationErr *errors.ValidationError

			if goerrors.As(err, &notFoundErr) {
				http.Error(w, "Student not found", http.StatusNotFound)
				return
			}
			if goerrors.As(err, &validationErr) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Gamer profile deleted successfully"))
	})
}
