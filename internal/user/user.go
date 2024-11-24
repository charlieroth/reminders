package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email string, passwordHash string, now time.Time) User {
	return User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now, UpdatedAt: now}
}

type CreateUserRequest struct {
	Email        string
	PasswordHash string
}

func NewCreateUserRequest(email string, passwordHash string) CreateUserRequest {
	return CreateUserRequest{Email: email, PasswordHash: passwordHash}
}

type UpdateUserRequest struct {
	Email        *string
	PasswordHash *string
}

func NewUpdateUserRequest(email *string, passwordHash *string) UpdateUserRequest {
	return UpdateUserRequest{Email: email, PasswordHash: passwordHash}
}
