package auth

import "context"

type AuthRepository interface {
	Store(ctx context.Context, app *Application) error
	FindKeyById(ctx context.Context, keyId string) (*Application, error)
	UpdateLastUsed(ctx context.Context, keyId string) error
}
