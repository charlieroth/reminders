package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/charlieroth/reminders/internal/ports"
)

type AuthService struct {
	repo ports.SessionRepository
}

func NewAuthService(repo ports.SessionRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Implements the AuthService.Login method
func (as *AuthService) Login(ctx context.Context, req domain.CreateSessionRequest) (domain.Session, error) {
	return as.repo.CreateSession(ctx, req)
}

// Implements the AuthService.Logout method
func (as *AuthService) Logout(ctx context.Context, req domain.DeleteSessionRequest) error {
	return as.repo.DeleteSession(ctx, req)
}

// Implements the AuthService.Refresh method
func (as *AuthService) Refresh(ctx context.Context, req domain.RefreshSessionRequest) (domain.Session, error) {
	return as.repo.RefreshSession(ctx, req)
}

// Implements the AuthService.GetSession method
func (as *AuthService) GetSession(ctx context.Context, req domain.GetSessionRequest) (domain.Session, error) {
	return as.repo.GetSession(ctx, req)
}

// Implements the AuthService.RevokeSession method
func (as *AuthService) RevokeSession(ctx context.Context, req domain.RevokeSessionRequest) error {
	return as.repo.RevokeSession(ctx, req)
}
