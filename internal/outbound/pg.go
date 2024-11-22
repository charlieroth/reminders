package outbound

import (
	"context"
	"database/sql"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/google/uuid"
)

type Pg struct {
	db *sql.DB
}

func NewPg(db *sql.DB) *Pg {
	return &Pg{db: db}
}

func (pg *Pg) SaveTask(ctx context.Context, task *task.Task) (uuid.UUID, error) {
	return uuid.New(), nil
}

// Implement the TaskRepository interface
func (pg *Pg) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	return task.Task{}, nil
}
