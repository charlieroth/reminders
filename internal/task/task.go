package task

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TaskTitle struct {
	Title string
}

func NewTaskTitle(title string) (TaskTitle, error) {
	trimmed := strings.Trim(title, " ")
	if trimmed == "" {
		return TaskTitle{}, &TaskTitleEmptyError{}
	}

	return TaskTitle{Title: trimmed}, nil
}

type Task struct {
	ID        uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(id uuid.UUID, title string, now time.Time) Task {
	return Task{ID: id, Title: title, CreatedAt: now, UpdatedAt: now}
}

type CreateTaskRequest struct {
	Title string
}

func NewCreateTaskRequest(title string) CreateTaskRequest {
	return CreateTaskRequest{Title: title}
}

type TaskTitleEmptyError struct{}

func (e *TaskTitleEmptyError) Error() string {
	return "title cannot be empty"
}
