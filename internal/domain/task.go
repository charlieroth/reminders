package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTaskRequest struct {
	Title string
}

type UpdateTaskRequest struct {
	Title     *string
	Completed *bool
}

type TaskTitleEmptyError struct{}

func (e *TaskTitleEmptyError) Error() string {
	return "title cannot be empty"
}
