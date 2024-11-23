package task

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

func NewTask(id uuid.UUID, title string, completed bool, now time.Time) Task {
	return Task{ID: id, Title: title, Completed: completed, CreatedAt: now, UpdatedAt: now}
}

type CreateTaskRequest struct {
	Title string
}

func NewCreateTaskRequest(title string) CreateTaskRequest {
	return CreateTaskRequest{Title: title}
}

type UpdateTaskRequest struct {
	Title     *string
	Completed *bool
}

func NewUpdateTaskRequest(title *string, completed *bool) UpdateTaskRequest {
	return UpdateTaskRequest{Title: title, Completed: completed}
}

type TaskTitleEmptyError struct{}

func (e *TaskTitleEmptyError) Error() string {
	return "title cannot be empty"
}
