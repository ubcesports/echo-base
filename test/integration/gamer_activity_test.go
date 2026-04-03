//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ubcesports/echo-base/internal/models"
)

func TestGamerActivityEndpoints(t *testing.T) {
	cleanupTestData(t)

	req := models.CreateGamerProfileRequest{
		StudentNumber:  "22222222",
		FirstName:      "Alice",
		LastName:       "Smith",
		MembershipTier: 1,
		Banned:         ptrBool(false),
	}
	rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to create test profile: %s", rr.Body.String())
	}

	t.Run("start activity", func(t *testing.T) {
		req := models.CreateActivityRequest{
			StudentNumber: "22222222",
			PCNumber:      1,
			Game:          "League of Legends",
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/activity", req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		var activity models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activity); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if activity.Game != "League of Legends" {
			t.Errorf("expected game League of Legends, got %s", activity.Game)
		}
		if activity.EndedAt != nil {
			t.Error("expected ended_at to be nil")
		}
	})

	t.Run("get active sessions", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/activity/all/get-active-pcs", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var activities []models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activities); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(activities) != 1 {
			t.Errorf("expected 1 active session, got %d", len(activities))
		}
	})

	t.Run("end activity", func(t *testing.T) {
		req := models.UpdateActivityRequest{
			PCNumber: 1,
			ExecName: "TestExec",
		}

		rr := makeRequest(t, http.MethodPatch, "/v1/api/activity/update/22222222", req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		var activity models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activity); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if activity.EndedAt == nil {
			t.Error("expected ended_at to be set")
		}
		if activity.ExecName == nil || *activity.ExecName != "TestExec" {
			t.Errorf("expected exec_name TestExec, got %v", activity.ExecName)
		}
	})

	t.Run("get exec leaderboard", func(t *testing.T) {
		profileReq := models.CreateGamerProfileRequest{
			StudentNumber:  "33333333",
			FirstName:      "Bob",
			LastName:       "Jones",
			MembershipTier: 1,
			Banned:         ptrBool(false),
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/gamer", profileReq)
		if rr.Code != http.StatusCreated {
			t.Fatalf("failed to create leaderboard test profile: %s", rr.Body.String())
		}

		req := models.CreateActivityRequest{
			StudentNumber: "33333333",
			PCNumber:      2,
			Game:          "VALORANT",
		}

		rr = makeRequest(t, http.MethodPost, "/v1/api/activity", req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		endReq := models.UpdateActivityRequest{
			PCNumber: 2,
			ExecName: "TestExec",
		}

		rr = makeRequest(t, http.MethodPatch, "/v1/api/activity/update/33333333", endReq)
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		req = models.CreateActivityRequest{
			StudentNumber: "33333333",
			PCNumber:      3,
			Game:          "Rocket League",
		}

		rr = makeRequest(t, http.MethodPost, "/v1/api/activity", req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		endReq = models.UpdateActivityRequest{
			PCNumber: 3,
			ExecName: "AnotherExec",
		}

		rr = makeRequest(t, http.MethodPatch, "/v1/api/activity/update/33333333", endReq)
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		rr = makeRequest(t, http.MethodGet, "/v1/api/activity/all/leaderboard", nil)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var leaderboard []models.ExecLeaderboardEntry
		if err := json.NewDecoder(rr.Body).Decode(&leaderboard); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(leaderboard) < 2 {
			t.Fatalf("expected at least 2 leaderboard entries, got %d", len(leaderboard))
		}

		if leaderboard[0].ExecName != "TestExec" || leaderboard[0].SignoutCount < 2 {
			t.Errorf("expected top leaderboard entry to be TestExec with count >= 2, got %+v", leaderboard[0])
		}
	})

	t.Run("get activities by student", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/activity/22222222", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var activities []models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activities); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(activities) != 1 {
			t.Errorf("expected 1 activity, got %d", len(activities))
		}
	})

	t.Run("get recent activities", func(t *testing.T) {
		rr := makeRequest(t, http.MethodGet, "/v1/api/activity/all/recent?page=1&limit=10", nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var activities []models.GamerActivity
		if err := json.NewDecoder(rr.Body).Decode(&activities); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(activities) == 0 {
			t.Error("expected at least one activity")
		}

		foundAlice := false
		for _, activity := range activities {
			if activity.FirstName != nil && *activity.FirstName == "Alice" {
				foundAlice = true
				break
			}
		}

		if !foundAlice {
			t.Error("expected first_name Alice to be included from JOIN")
		}
	})

	t.Run("start activity for non-existent student", func(t *testing.T) {
		req := models.CreateActivityRequest{
			StudentNumber: "99999999",
			PCNumber:      1,
			Game:          "Test",
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/activity", req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("invalid student number", func(t *testing.T) {
		req := models.CreateActivityRequest{
			StudentNumber: "123",
			PCNumber:      1,
			Game:          "Test",
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/activity", req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("missing game", func(t *testing.T) {
		req := models.CreateActivityRequest{
			StudentNumber: "22222222",
			PCNumber:      1,
		}

		rr := makeRequest(t, http.MethodPost, "/v1/api/activity", req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
