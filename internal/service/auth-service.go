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
func (as *AuthService) Logout(ctx context.Context, req session.DeleteSessionRequest) error {
	return as.repo.DeleteSession(ctx, req)
}

// Implements the AuthService.Refresh method
func (as *AuthService) Refresh(ctx context.Context, req session.RefreshSessionRequest) (session.Session, error) {
	return as.repo.RefreshSession(ctx, req)
}

// Implements the AuthService.GetSession method
func (as *AuthService) GetSession(ctx context.Context, req session.GetSessionRequest) (session.Session, error) {
	return as.repo.GetSession(ctx, req)
}

// Implements the AuthService.RevokeSession method
func (as *AuthService) RevokeSession(ctx context.Context, req session.RevokeSessionRequest) error {
	return as.repo.RevokeSession(ctx, req)
}
