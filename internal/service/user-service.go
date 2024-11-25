package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/ports"
	"github.com/charlieroth/reminders/internal/user"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository ports.UserRepository
}

func NewUserService(userRepository ports.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	return us.userRepository.GetUserByEmail(ctx, email)
}

func (us *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	return us.userRepository.GetUserByID(ctx, id)
}

func (us *UserService) CreateUser(ctx context.Context, req user.CreateUserRequest) (user.User, error) {
	return us.userRepository.CreateUser(ctx, req)
}

func (us *UserService) GetUser(ctx context.Context, id uuid.UUID) (user.User, error) {
	return us.userRepository.GetUser(ctx, id)
}

func (us *UserService) GetUsers(ctx context.Context) ([]user.User, error) {
	return us.userRepository.GetUsers(ctx)
}

func (us *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req user.UpdateUserRequest) (user.User, error) {
	return us.userRepository.UpdateUser(ctx, id, req)
}
