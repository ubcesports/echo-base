//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ubcesports/echo-base/internal/models"
)

func TestPagination(t *testing.T) {
	cleanupTestData(t)

	req := models.CreateGamerProfileRequest{
		StudentNumber:  "44444444",
		FirstName:      "Test",
		LastName:       "User",
		MembershipTier: 0,
		Banned:         ptrBool(false),
	}
	makeRequest(t, http.MethodPost, "/v1/api/gamer", req)

	for i := 1; i <= 15; i++ {
		actReq := models.CreateActivityRequest{
			StudentNumber: "44444444",
			PCNumber:      i,
			Game:          fmt.Sprintf("Game%d", i),
		}
		makeRequest(t, http.MethodPost, "/v1/api/activity", actReq)
	}

	t.Run("pagination limit works", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/activity/all/recent?page=1&limit=5", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var activities []models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activities); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(activities) != 5 {
			t.Errorf("expected 5 activities, got %d", len(activities))
		}
	})

	t.Run("invalid pagination", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/activity/all/recent?page=0&limit=10", nil)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
