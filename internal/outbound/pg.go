package outbound

import (
	"context"
	"database/sql"
	"time"

	"github.com/charlieroth/reminders/internal/task"
	"github.com/google/uuid"
)

type Pg struct {
	db *sql.DB
}

func NewPg(db *sql.DB) *Pg {
	return &Pg{db: db}
}

// Implements the TaskRepository.CreateTask method
func (pg *Pg) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	taskId := uuid.New()
	row := pg.db.QueryRowContext(ctx, `
		INSERT INTO tasks (id, title, completed)
		VALUES ($1, $2, $3)
		RETURNING id, title, created_at, updated_at, completed
	`, taskId, req.Title, false)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt, &t.Completed); err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskRepository.GetTask method
func (pg *Pg) GetTask(ctx context.Context, id uuid.UUID) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		SELECT id, title, completed, created_at, updated_at FROM tasks WHERE id = $1
	`, id)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskRepository.ListTasks method
func (pg *Pg) ListTasks(ctx context.Context) ([]task.Task, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, title, completed, created_at, updated_at FROM tasks
	`)
	if err != nil {
		return []task.Task{}, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return []task.Task{}, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Implements the TaskRepository.UpdateTask method
func (pg *Pg) UpdateTask(ctx context.Context, id uuid.UUID, req task.UpdateTaskRequest) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		UPDATE tasks 
		SET 
			title = COALESCE($1, title),
			completed = COALESCE($2, completed),
			updated_at = $3 
		WHERE id = $4
		RETURNING id, title, created_at, updated_at, completed
	`, req.Title, req.Completed, time.Now().UTC(), id)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt, &t.Completed); err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}
