package ports

import (
	"context"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (task.Task, error)
	ListTasks(ctx context.Context) ([]task.Task, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (task.Task, error)
	ListTasks(ctx context.Context) ([]task.Task, error)
}
