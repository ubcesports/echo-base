//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ubcesports/echo-base/internal/models"
)

func TestGamerProfileEndpoints(t *testing.T) {
	cleanupTestData(t)

	t.Run("create gamer profile", func(t *testing.T) {
		req := models.CreateGamerProfileRequest{
			StudentNumber:  "12345678",
			FirstName:      "John",
			LastName:       "Doe",
			MembershipTier: 1,
			Banned:         ptrBool(false),
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		var profile models.GamerProfile
		if err := json.NewDecoder(rr.Body).Decode(&profile); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if profile.StudentNumber != "12345678" {
			t.Errorf("expected student_number 12345678, got %s", profile.StudentNumber)
		}

		if profile.MembershipExpiryDate == nil {
			t.Error("expected membership_expiry_date to be set")
		} else {
			expectedYear := time.Now().Year()
			if time.Now().Month() >= time.May {
				expectedYear++
			}
			if profile.MembershipExpiryDate.Year() != expectedYear {
				t.Errorf("expected expiry year %d, got %d", expectedYear, profile.MembershipExpiryDate.Year())
			}
			if profile.MembershipExpiryDate.Month() != time.May || profile.MembershipExpiryDate.Day() != 1 {
				t.Errorf("expected expiry date May 1st, got %s", profile.MembershipExpiryDate.Format("2006-01-02"))
			}
		}
	})

	t.Run("get gamer profile", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/gamer/12345678", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var profile models.GamerProfile
		if err := json.NewDecoder(rr.Body).Decode(&profile); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if profile.FirstName != "John" {
			t.Errorf("expected first_name John, got %s", profile.FirstName)
		}
	})

	t.Run("get non-existent profile", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/gamer/99999999", nil)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("update existing profile", func(t *testing.T) {
		req := models.CreateGamerProfileRequest{
			StudentNumber:  "12345678",
			FirstName:      "Jane",
			LastName:       "Doe",
			MembershipTier: 2,
			Banned:         ptrBool(false),
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
		}

		var profile models.GamerProfile
		if err := json.NewDecoder(rr.Body).Decode(&profile); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if profile.FirstName != "Jane" {
			t.Errorf("expected first_name Jane, got %s", profile.FirstName)
		}
		if profile.MembershipTier != 2 {
			t.Errorf("expected membership_tier 2, got %d", profile.MembershipTier)
		}
	})

	t.Run("invalid student number", func(t *testing.T) {
		req := models.CreateGamerProfileRequest{
			StudentNumber:  "123",
			FirstName:      "Test",
			LastName:       "User",
			MembershipTier: 1,
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("missing first name", func(t *testing.T) {
		req := models.CreateGamerProfileRequest{
			StudentNumber:  "11111111",
			LastName:       "User",
			MembershipTier: 1,
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("delete gamer profile", func(t *testing.T) {
		rr := makeRequest(t, http.MethodDelete, "/v1/api/gamer/12345678", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		rr = makeRequest(t, http.MethodGet, "/v1/api/gamer/12345678", nil)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d after deletion, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("delete non-existent profile", func(t *testing.T) {
		rr := makeRequest(t, http.MethodDelete, "/v1/api/gamer/99999999", nil)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}
