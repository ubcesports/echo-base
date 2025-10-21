package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

const (
	KeyIDLength      = 6
	SecretLength     = 32
	APIKeyPrefix     = "api_"
	MaxAppNameLength = 100
)

var (
	validAppNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type AuthService interface {
	GenerateAPIKey(ctx context.Context, appName string) (*auth.APIKey, error)
	ValidateAPIKey(ctx context.Context, apiKey string) (string, error)
}

type authService struct {
	repo auth.AuthRepository
}

func NewAuthService(repo auth.AuthRepository) *authService {
	return &authService{repo: repo}
}

func (s *authService) GenerateAPIKey(ctx context.Context, appName string) (*auth.APIKey, error) {
	if err := s.validateAppName(appName); err != nil {
		return nil, err
	}

	keyId, secret, err := s.genereateCredentials()
	if err != nil {
		return nil, err
	}

	hashedSecret := s.hashSecret(secret)
	app := &auth.Application{
		AppName:   appName,
		KeyId:     keyId,
		HashedKey: hashedSecret,
	}

	if err := s.repo.Store(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}
	return &auth.APIKey{
		KeyId:   keyId,
		APIKey:  fmt.Sprintf("api_%s.%s", keyId, secret),
		AppName: appName,
	}, nil
}

func (s *authService) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {

	keyId, secret, err := s.parseAPIKey(apiKey)

	app, err := s.repo.FindKeyById(ctx, keyId)
	if err != nil {
		return "", err
	}
	if !s.verifySecret(secret, app.HashedKey) {
		return "", fmt.Errorf("invalid api key")
	}

	go func() {
		s.repo.UpdateLastUsed(context.Background(), keyId)
	}()

	return app.AppName, nil
}

func (s *authService) genereateCredentials() (string, string, error) {
	keyIDBytes := make([]byte, KeyIDLength)
	_, err := rand.Read(keyIDBytes)
	if err != nil {
		return "", "", err
	}
	keyId := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(keyIDBytes))

	secretBytes := make([]byte, SecretLength)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return "", "", err
	}
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)
	return keyId, secret, nil
}

func (s *authService) hashSecret(secret string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedSecret := hasher.Sum(nil)
	return hashedSecret
}

func (s *authService) validateAppName(appName string) error {
	if appName == "" {
		return fmt.Errorf("app_name is required")
	}

	if len(appName) > MaxAppNameLength {
		return fmt.Errorf("app_name must be %d characters or less", MaxAppNameLength)
	}

	if !validAppNameRegex.MatchString(appName) {
		return fmt.Errorf("app_name can only contain letters, numbers, hyphens, and underscores")
	}

	return nil
}

func (s *authService) parseAPIKey(apiKey string) (string, string, error) {
	if !strings.HasPrefix(apiKey, "api_") {
		return "", "", fmt.Errorf("invalid api_key format, missing \"api_\" prefix")
	}

	keyParts := strings.Split(strings.TrimPrefix(apiKey, "api_"), ".")
	if len(keyParts) != 2 {
		return "", "", fmt.Errorf("invalid api_key format, missing parts")
	}

	keyId := keyParts[0]
	secret := keyParts[1]
	return keyId, secret, nil
}

func (s *authService) verifySecret(secret string, hashedSecret []byte) bool {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	actualHash := hasher.Sum(nil)

	return subtle.ConstantTimeCompare(hashedSecret, actualHash) == 1
}
