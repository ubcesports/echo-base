package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

type mockAuthService struct {
	validKeys map[string]string // apiKey -> appName
}

func (m *mockAuthService) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {
	if appName, exists := m.validKeys[apiKey]; exists {
		return appName, nil
	}
	return "", fmt.Errorf("invalid API key")
}

func (m *mockAuthService) GenerateAPIKey(ctx context.Context, appName string) (*auth.APIKey, error) {
	return nil, nil
}

func TestAuthMiddleware(t *testing.T) {
	mockService := &mockAuthService{
		validKeys: map[string]string{
			"valid-api-key": "test-app",
		},
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appName := r.Context().Value("appName").(string)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success: " + appName))
	})

	testCases := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid auth header",
			authHeader:     "Basic token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Bearer token",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Bearer token",
			authHeader:     "Bearer valid-api-key",
			expectedStatus: http.StatusOK,
			expectedBody:   "success: test-app",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := AuthMiddleware(testHandler, mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
