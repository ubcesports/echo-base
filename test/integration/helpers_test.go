//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/services"
)

var (
	testServer http.Handler
	testAPIKey string
)

func TestMain(m *testing.M) {
	config.LoadEnv("../../.env.test")
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

func ptrBool(b bool) *bool {
	return &b
}
