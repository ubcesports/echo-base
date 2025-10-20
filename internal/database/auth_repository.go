package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ubcesports/echo-base/internal/interfaces/auth"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Store(ctx context.Context, app *auth.Application) error {
	query := `
        INSERT INTO application (app_name, key_id, hashed_key)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(ctx, query,
		app.AppName, app.KeyId, app.HashedKey)
	return err
}

func (r *AuthRepository) FindKeyById(ctx context.Context, keyId string) (*auth.Application, error) {
	query := `
        SELECT app_name, key_id, hashed_key
        FROM application 
        WHERE key_id = $1
    `

	app := &auth.Application{}

	err := r.db.QueryRowContext(ctx, query, keyId).Scan(
		&app.KeyId, &app.AppName, &app.HashedKey,
	)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AuthRepository) UpdateLastUsed(ctx context.Context, keyId string) error {
	query := `UPDATE application SET last_used_at = $1 WHERE key_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), keyId)
	return err
}
