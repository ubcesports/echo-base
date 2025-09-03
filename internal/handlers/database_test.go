// filepath: internal/handlers/database_test.go
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/echo-base/internal/database"
)

func TestDatabasePing(t *testing.T) {
	database.Init()
	defer database.Close()

	req, err := http.NewRequest("GET", "/db/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DatabasePing)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("Expected Content-Type %s, got %s", expected, ct)
	}

	var response DatabasePingResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.ResponseTime == "" {
		t.Error("Response time should not be empty")
	}
}
