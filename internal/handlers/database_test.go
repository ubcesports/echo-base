package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/echo-base/internal/tests"
)

func TestDatabasePing(t *testing.T) {
	tests.SetupTestDB(t)

	req := tests.CreateTestRequest(t, "GET", "/db/ping", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DatabasePing)
	handler.ServeHTTP(rr, req)

	var response DatabasePingResponse
	tests.AssertResponse(t, rr, http.StatusOK, &response)

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.ResponseTime == "" {
		t.Error("Response time should not be empty")
	}
}
