package domain

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewList(id uuid.UUID, name string, now time.Time) List {
	return List{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type CreateListRequest struct {
	Name string
}

func NewCreateListRequest(name string) CreateListRequest {
	return CreateListRequest{
		Name: name,
	}
}

type UpdateListRequest struct {
	Name string
}
