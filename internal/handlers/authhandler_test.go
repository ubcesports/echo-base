package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

// MockAuthService implements services.AuthService for testing
type MockAuthService struct {
	GenerateAPIKeyFunc func(ctx context.Context, appName string) (*auth.APIKey, error)
}

func (m *MockAuthService) GenerateAPIKey(ctx context.Context, appName string) (*auth.APIKey, error) {
	return m.GenerateAPIKeyFunc(ctx, appName)
}

func (m *MockAuthService) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {
	return "", nil
}

func TestGenerateAPIKeyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		mockReturn     *auth.APIKey
		mockErr        error
		wantStatus     int
		wantBodySubstr string
	}{
		{
			name:           "valid POST",
			method:         http.MethodPost,
			body:           GenerateKeyRequest{AppName: "testapp"},
			mockReturn:     &auth.APIKey{KeyId: "abc123", APIKey: "api_abc123.secret", AppName: "testapp"},
			wantStatus:     http.StatusOK,
			wantBodySubstr: `"api_abc123.secret"`,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			body:           nil,
			wantStatus:     http.StatusMethodNotAllowed,
			wantBodySubstr: "Method not allowed",
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			body:           "{invalid-json}",
			wantStatus:     http.StatusBadRequest,
			wantBodySubstr: "Invalid request body",
		},
		{
			name:           "service error",
			method:         http.MethodPost,
			body:           GenerateKeyRequest{AppName: "testapp"},
			mockErr:        errors.New("service error"),
			wantStatus:     http.StatusInternalServerError,
			wantBodySubstr: "Error generating API key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{
				GenerateAPIKeyFunc: func(ctx context.Context, appName string) (*auth.APIKey, error) {
					return tt.mockReturn, tt.mockErr
				},
			}

			var reqBody []byte
			switch v := tt.body.(type) {
			case GenerateKeyRequest:
				reqBody, _ = json.Marshal(v)
			case string:
				reqBody = []byte(v)
			case nil:
				reqBody = nil
			}

			req := httptest.NewRequest(tt.method, "/api-key", bytes.NewReader(reqBody))
			rec := httptest.NewRecorder()

			handler := GenerateAPIKey(mockService)
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBodySubstr != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tt.wantBodySubstr)) {
				t.Errorf("response body %q does not contain %q", rec.Body.String(), tt.wantBodySubstr)
			}
		})
	}
}
