package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/echo-base/internal/tests"
)

func TestHealthCheck(t *testing.T) {
	tests.SetupTestDB(t)

	req := tests.CreateTestRequest(t, "GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	var response HealthResponse
	tests.AssertResponse(t, rr, http.StatusOK, &response)

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.Database != "ok" {
		t.Errorf("Expected database status 'ok', got '%s'", response.Database)
	}
}
