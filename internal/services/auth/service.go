package auth

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateAPI(ctx context.Context, appName string) (*APIKey, error) {
	return nil, nil
}

func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {
	return "", nil
}
