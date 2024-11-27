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

type CreateListRequest struct {
	Name string
}

type UpdateListRequest struct {
	Name string
}
