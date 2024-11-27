package domain

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

type UserList struct {
	UserID uuid.UUID
	ListID uuid.UUID
}

type UserTask struct {
	UserID uuid.UUID
	TaskID uuid.UUID
}

type CreateUserRequest struct {
	Email        string
	PasswordHash string
}

type UpdateUserRequest struct {
	Email        *string
	PasswordHash *string
}
