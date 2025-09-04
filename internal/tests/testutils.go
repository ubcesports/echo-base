package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ubcesports/echo-base/internal/database"
)

func SetupTestDB(t *testing.T) {
	os.Setenv("EB_DSN", "postgresql://user:pass@localhost:5433/echobase_test?sslmode=disable")

	database.Init()
	t.Cleanup(func() {
		database.Close()
	})
}

func CreateTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, body interface{}) {
	if status := rr.Code; status != expectedStatus {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	if body != nil {
		if err := json.NewDecoder(rr.Body).Decode(body); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}
	}
}
