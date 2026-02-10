//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/models"
)

func TestExpiredMembership(t *testing.T) {
	cleanupTestData(t)

	req := models.CreateGamerProfileRequest{
		StudentNumber:  "33333333",
		FirstName:      "Bob",
		LastName:       "Jones",
		MembershipTier: 2,
		Banned:         ptrBool(false),
	}
	rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to create test profile: %s", rr.Body.String())
	}

	_, err := database.DB.Exec("UPDATE gamer_profile SET membership_expiry_date = $1 WHERE student_number = $2",
		time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), "33333333")
	if err != nil {
		t.Fatalf("failed to update expiry date: %v", err)
	}

	t.Run("expired membership blocks check-in", func(t *testing.T) {
		req := models.CreateActivityRequest{
			StudentNumber: "33333333",
			PCNumber:      1,
			Game:          "Test",
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/activity", req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d: %s", http.StatusForbidden, rr.Code, rr.Body.String())
		}

		body := rr.Body.String()
		if body == "" || len(body) < 10 {
			t.Error("expected detailed error message for expired membership")
		}
	})
}
