package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Token       string
	CreatedAt   time.Time
	RefreshedAt time.Time
	ExpiresAt   time.Time
	UserAgent   string
	Active      bool
}

func NewSession(userID uuid.UUID, token string, userAgent string) *Session {
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		UserAgent: userAgent,
	}
}

type CreateSessionRequest struct {
	UserID    uuid.UUID
	Token     string
	UserAgent string
}

type RefreshSessionRequest struct {
	UserID    uuid.UUID
	Token     string
	UserAgent string
}

type InvalidateSessionRequest struct {
	UserID uuid.UUID
	Token  string
}
