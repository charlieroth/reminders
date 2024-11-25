package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/ports"
	"github.com/charlieroth/reminders/internal/session"
)

type AuthService struct {
	repo ports.SessionRepository
}

func NewAuthService(repo ports.SessionRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Implements the AuthService.Login method
func (as *AuthService) Login(ctx context.Context, req session.CreateSessionRequest) (session.Session, error) {
	return as.repo.CreateSession(ctx, req)
}

// Implements the AuthService.Logout method
func (as *AuthService) Logout(ctx context.Context, req session.InvalidateSessionRequest) error {
	return as.repo.InvalidateSession(ctx, req)
}

// Implements the AuthService.LogoutByEmail method
func (as *AuthService) LogoutByEmail(ctx context.Context, email string) error {
	return as.repo.InvalidateSessionByEmail(ctx, session.InvalidateSessionByEmailRequest{Email: email})
}

// Implements the AuthService.Refresh method
func (as *AuthService) Refresh(ctx context.Context, req session.RefreshSessionRequest) (session.Session, error) {
	return as.repo.RefreshSession(ctx, req)
}

// Implements the AuthService.GetSessions method
func (as *AuthService) GetSessions(ctx context.Context) ([]session.Session, error) {
	return as.repo.GetSessions(ctx)
}
