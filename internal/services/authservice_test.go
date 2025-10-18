package services

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

type mockAuthRepository struct {
	applications map[string]*auth.Application
}

func (m *mockAuthRepository) Store(ctx context.Context, app *auth.Application) error {
	m.applications[app.KeyId] = app
	return nil
}
func (m *mockAuthRepository) FindKeyById(ctx context.Context, keyId string) (*auth.Application, error) {
	app, exists := m.applications[keyId]
	if !exists {
		return nil, fmt.Errorf("error")
	}
	return app, nil
}
func (m *mockAuthRepository) UpdateLastUsed(ctx context.Context, keyId string) error {
	return nil
}

func TestGenerateAPIKey(t *testing.T) {
	mockRepo := &mockAuthRepository{
		applications: make(map[string]*auth.Application),
	}
	authService := NewAuthService(mockRepo)

	// Valid Key
	apiKey, err := authService.GenerateAPIKey(context.Background(), "test-app")
	assert.NoError(t, err)
	assert.Equal(t, "test-app", apiKey.AppName)

	// Invalid key, no app name
	apiKey, err = authService.GenerateAPIKey(context.Background(), "")
	assert.Error(t, err)

	// Invalid key, too long
	apiKey, err = authService.GenerateAPIKey(context.Background(), strings.Repeat("a", 101))
	assert.Error(t, err)
}

func TestValidateAPIKey(t *testing.T) {
	mockRepo := &mockAuthRepository{
		applications: make(map[string]*auth.Application),
	}

	authService := NewAuthService(mockRepo)
	rawSecret := "supersecret"
	hashedSecret := authService.hashSecret(rawSecret)
	mockApp := auth.Application{
		KeyId:     "testid",
		HashedKey: hashedSecret,
		AppName:   "test-app",
	}
	mockRepo.Store(context.Background(), &mockApp)

	// Valid API key
	apiKey := fmt.Sprintf("api_%s.%s", mockApp.KeyId, rawSecret)
	gotAppName, err := authService.ValidateAPIKey(context.Background(), apiKey)
	assert.NoError(t, err)
	assert.Equal(t, mockApp.AppName, gotAppName)

	// Invalid secret
	badApiKey := fmt.Sprintf("api_%s.%s", mockApp.KeyId, "wrongsecret")
	_, err = authService.ValidateAPIKey(context.Background(), badApiKey)
	assert.Error(t, err)
}
