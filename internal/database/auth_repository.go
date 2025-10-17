package database

import (
	"context"
	"database/sql"

	"github.com/ubcesports/echo-base/internal/services/auth"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Store(ctx context.Context, app *auth.Application) error {
	return nil
}

func (r *AuthRepository) FindKeyById(ctx context.Context, keyId string) (*auth.Application, error) {
	return nil, nil
}

func (r *AuthRepository) UpdateLastUsed(ctx context.Context, keyId string) error {
	return nil
}
