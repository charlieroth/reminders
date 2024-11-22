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

func (pg *Pg) InsertTask(ctx context.Context, title string) (uuid.UUID, error) {
	id := uuid.New()
	row := pg.db.QueryRowContext(ctx, `
		INSERT INTO tasks (id, title)
		VALUES ($1, $2)
		RETURNING id
	`, id, title)

	var createdID uuid.UUID
	if err := row.Scan(&createdID); err != nil {
		return uuid.UUID{}, err
	}
	return createdID, nil
}

func (pg *Pg) SelectTaskByID(ctx context.Context, id uuid.UUID) (task.Task, error) {
	row := pg.db.QueryRowContext(ctx, `
		SELECT id, title, created_at, updated_at FROM tasks WHERE id = $1
	`, id)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return task.Task{}, err
	}
	return t, nil
}

func (pg *Pg) SelectTasks(ctx context.Context) ([]task.Task, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, title, created_at, updated_at FROM tasks
	`)
	if err != nil {
		return []task.Task{}, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return []task.Task{}, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Implements the TaskRepository.CreateTask method
func (pg *Pg) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	taskId, err := pg.InsertTask(ctx, req.Title)
	if err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	task := task.NewTask(taskId, req.Title, time.Now().UTC())
	return task, nil
}

// Implements the TaskRepository.GetTask method
func (pg *Pg) GetTask(ctx context.Context, id uuid.UUID) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	t, err := pg.SelectTaskByID(ctx, id)
	if err != nil {
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
		SELECT id, title, created_at, updated_at FROM tasks
	`)
	if err != nil {
		return []task.Task{}, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return []task.Task{}, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
