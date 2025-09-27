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

// Does not include auto cleanup
func SetupTestDB() {
	os.Setenv("EB_DSN", "postgresql://user:pass@localhost:5433/echobase_test?sslmode=disable")
	database.Init()
}

func SetupTestDBForTest(t *testing.T) {
	SetupTestDB()
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

	req := httptest.NewRequest(method, url, &reqBody)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, body interface{}) {
	if status := rr.Code; status != expectedStatus {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	if body != nil {
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", ct)
		}

		if err := json.NewDecoder(rr.Body).Decode(body); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}
	}
}

func ExecuteTestRequest(t *testing.T, router http.Handler, method, url string, body interface{}) *httptest.ResponseRecorder {
	req := CreateTestRequest(t, method, url, body)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
