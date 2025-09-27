package handlers

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/ubcesports/echo-base/internal/tests"
)

func TestGenerateAPIKey(t *testing.T) {
	tests.SetupTestDBForTest(t)

	testCases := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		expectKey      bool
	}{
		{
			name:           "Valid request",
			method:         "POST",
			body:           map[string]string{"app_name": "test-app"},
			expectedStatus: http.StatusOK,
			expectKey:      true,
		},
		{
			name:           "Invalid method",
			method:         "GET",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectKey:      false,
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			body:           "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectKey:      false,
		},
		{
			name:           "Invalid app name",
			method:         "POST",
			body:           map[string]string{"app_name": ""},
			expectedStatus: http.StatusBadRequest,
			expectKey:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body == "invalid-json" {
				req, _ = http.NewRequest(tc.method, "/admin/generate-key", strings.NewReader(`{"app_name": }`))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = tests.CreateTestRequest(t, tc.method, "/admin/generate-key", tc.body)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GenerateAPIKey)
			handler.ServeHTTP(rr, req)

			if tc.expectKey {
				var response GenerateKeyResponse
				tests.AssertResponse(t, rr, tc.expectedStatus, &response)
				if response.KeyID == "" {
					t.Error("KeyID should not be empty")
				}
				if !strings.HasPrefix(response.APIKey, "api_") {
					t.Error("API key should start with 'api_'")
				}
				if response.AppName != "test-app" {
					t.Errorf("Expected app name 'test-app', got '%s'", response.AppName)
				}
			} else {
				if rr.Code != tc.expectedStatus {
					t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
				}
			}
		})
	}
}

/*
for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var response GenerateKeyResponse

			if tc.rawBody != "" {
				req, _ := http.NewRequest(tc.method, "/admin/generate-key", strings.NewReader(tc.rawBody))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(GenerateAPIKey)
				handler.ServeHTTP(rr, req)
				tests.AssertResponse(t, rr, tc.expectedStatus, nil)
			} else {
				tests.ExecuteTestRequest(t, http.HandlerFunc(GenerateAPIKey), tc.method, "/admin/generate-key", tc.body, tc.expectedStatus, &response)
			}

			if tc.expectKey {
				if response.KeyID == "" {
					t.Error("KeyID should not be empty")
				}
				if !strings.HasPrefix(response.APIKey, "api_") {
					t.Error("API key should start with 'api_'")
				}
				if response.AppName != "test-app" {
					t.Errorf("Expected app name 'test-app', got '%s'", response.AppName)
				}
			}
		})
	}
}
*/
