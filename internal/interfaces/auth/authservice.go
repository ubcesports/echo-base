package auth

import "context"

type AuthService interface {
	GenerateAPIKey(ctx context.Context, appName string) (*APIKey, error)
	ValidateAPIKey(ctx context.Context, apiKey string) (string, error)
}
