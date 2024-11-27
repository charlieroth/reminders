package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/charlieroth/reminders/internal/ports"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository ports.UserRepository
}

func NewUserService(userRepository ports.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return us.userRepository.GetUserByEmail(ctx, email)
}

func (us *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return us.userRepository.GetUserByID(ctx, id)
}

func (us *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.User, error) {
	return us.userRepository.CreateUser(ctx, req)
}

func (us *UserService) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return us.userRepository.GetUser(ctx, id)
}

func (us *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	return us.userRepository.GetUsers(ctx)
}

func (us *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (domain.User, error) {
	return us.userRepository.UpdateUser(ctx, id, req)
}
