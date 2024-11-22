package ports

import (
	"context"

	"github.com/charlieroth/reminders/internal/task"
)

type TaskService interface {
	CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error)
}
