package auth

import (
	"context"

	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

type Service struct {
	repo auth.AuthRepository
}

func NewService(repo auth.AuthRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateAPIKey(ctx context.Context, appName string) (*auth.APIKey, error) {
	return nil, nil
}

func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {
	return "", nil
}
