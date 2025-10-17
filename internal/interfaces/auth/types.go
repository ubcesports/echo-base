package auth

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	Id         uuid.UUID
	AppName    string
	KeyId      string
	HashedKey  []byte
	CreatedAt  time.Time
	LastUsedAt time.Time
}

type APIKey struct {
	KeyId   string
	APIKey  string
	AppName string
}
