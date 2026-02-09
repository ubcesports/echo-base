package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/models"
	"github.com/ubcesports/echo-base/internal/services"
)

var (
	testServer http.Handler
	testAPIKey string
)

func TestMain(m *testing.M) {
	config.LoadEnv("../.env.test")
	database.Init()
	defer database.Close()

	authRepo := database.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo)
	gamerProfileRepo := database.NewGamerProfileRepository(database.DB)
	gamerActivityRepo := database.NewGamerActivityRepository(database.DB)
	gamerProfileService := services.NewGamerProfileService(gamerProfileRepo)
	gamerActivityService := services.NewGamerActivityService(gamerActivityRepo, gamerProfileRepo)

	testServer = internal.NewServer(authService, gamerProfileService, gamerActivityService)

	apiKey, err := authService.GenerateAPIKey(context.Background(), "integration-test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate API key: %v\n", err)
		os.Exit(1)
	}
	testAPIKey = apiKey.APIKey

	code := m.Run()
	os.Exit(code)
}

func makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+testAPIKey)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testServer.ServeHTTP(rr, req)
	return rr
}

func cleanupTestData(t *testing.T) {
	_, err := database.DB.Exec("DELETE FROM gamer_activity")
	if err != nil {
		t.Logf("Warning: failed to clean gamer_activity: %v", err)
	}
	_, err = database.DB.Exec("DELETE FROM gamer_profile")
	if err != nil {
		t.Logf("Warning: failed to clean gamer_profile: %v", err)
	}
}

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

		if activities[0].FirstName == nil || *activities[0].FirstName != "Alice" {
			t.Error("expected first_name to be included from JOIN")
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

func ptrBool(b bool) *bool {
	return &b
}
