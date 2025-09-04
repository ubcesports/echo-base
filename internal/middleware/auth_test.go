package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/echo-base/internal/tests"
)

func TestAuthMiddleware(t *testing.T) {
	tests.SetupTestDB(t)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	testCases := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{"No auth header", "", http.StatusUnauthorized},
		{"Invalid auth header", "Basic token", http.StatusUnauthorized},
		{"Invalid Bearer token", "Bearer invalid", http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tests.CreateTestRequest(t, "GET", "/test", nil)

			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := AuthMiddleware(testHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}
