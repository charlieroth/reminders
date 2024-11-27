package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func NewSession(email string, refreshToken string) *Session {
	return &Session{
		ID:           uuid.New(),
		Email:        email,
		RefreshToken: refreshToken,
	}
}

type CreateSessionRequest struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	IsRevoked    bool
	ExpiresAt    time.Time
}

type RefreshSessionRequest struct {
	Email        string
	RefreshToken string
}

type RevokeSessionRequest struct {
	ID uuid.UUID
}

type DeleteSessionRequest struct {
	ID uuid.UUID
}

type GetSessionRequest struct {
	ID uuid.UUID
}
