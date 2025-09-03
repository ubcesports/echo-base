package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ubcesports/echo-base/internal/database"
)

func TestValidateAppName(t *testing.T) {
	tests := []struct {
		name     string
		appName  string
		expected bool
	}{
		{"Valid app name", "my-app", false},
		{"Valid with numbers", "app123", false},
		{"Valid with underscores", "my_app", false},
		{"Empty app name", "", true},
		{"Too long app name", strings.Repeat("a", MaxAppNameLength+1), true},
		{"Invalid characters", "my app!", true},
		{"Valid edge case", strings.Repeat("a", MaxAppNameLength), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAppName(test.appName)
			if (err != nil) != test.expected {
				t.Errorf("validateAppName() error = %v, wantErr %v", err, test.expected)
			}
		})
	}
}

func TestGenerateAPIKey(t *testing.T) {
	database.Init()
	defer database.Close()

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			method:         "POST",
			body:           `{"app_name": "test-app"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         "GET",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			body:           `{"app_name": }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid app name",
			method:         "POST",
			body:           `{"app_name": ""}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/admin/generate-key", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GenerateAPIKey)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response GenerateKeyResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

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
