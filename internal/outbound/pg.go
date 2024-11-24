package outbound

import (
	"context"
	"database/sql"
	"time"

	"github.com/charlieroth/reminders/internal/list"
	"github.com/charlieroth/reminders/internal/session"
	"github.com/charlieroth/reminders/internal/task"
	"github.com/charlieroth/reminders/internal/user"
	"github.com/google/uuid"
)

type Pg struct {
	db *sql.DB
}

func NewPg(db *sql.DB) *Pg {
	return &Pg{db: db}
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
// Implements the DatabaseService.StatusCheck method
func (pg *Pg) StatusCheck(ctx context.Context) error {
	// If the user doesn't give us a deadline set 1 second.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	for attempts := 1; ; attempts++ {
		if err := pg.db.Ping(); err == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	const q = `SELECT TRUE`
	var tmp bool
	return pg.db.QueryRowContext(ctx, q).Scan(&tmp)
}

// Implements the SessionRepository.CreateSession method
func (pg *Pg) CreateSession(ctx context.Context, req session.CreateSessionRequest) (session.Session, error) {
	return session.Session{}, nil
}

// Implements the SessionRepository.RefreshSession method
func (pg *Pg) RefreshSession(ctx context.Context, req session.RefreshSessionRequest) (session.Session, error) {
	return session.Session{}, nil
}

// Implements the SessionRepository.InvalidateSession method
func (pg *Pg) InvalidateSession(ctx context.Context, req session.InvalidateSessionRequest) (session.Session, error) {
	return session.Session{}, nil
}

// Implements the SessionRepository.GetSessions method
func (pg *Pg) GetSessions(ctx context.Context) ([]session.Session, error) {
	return []session.Session{}, nil
}

// Implements the UserRepository.CreateUser method
func (pg *Pg) CreateUser(ctx context.Context, req user.CreateUserRequest) (user.User, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return user.User{}, err
	}
	defer tx.Rollback()

	userId := uuid.New()
	row := pg.db.QueryRowContext(ctx, `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, created_at, updated_at
	`, userId, req.Email, req.PasswordHash)

	var u user.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return user.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Implements the UserRepository.GetUser method
func (pg *Pg) GetUser(ctx context.Context, id uuid.UUID) (user.User, error) {
	row := pg.db.QueryRowContext(ctx, `
		SELECT id, email, created_at, updated_at FROM users WHERE id = $1
	`, id)

	var u user.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Implements the UserRepository.GetUsers method
func (pg *Pg) GetUsers(ctx context.Context) ([]user.User, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, email, created_at, updated_at FROM users
	`)
	if err != nil {
		return []user.User{}, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return []user.User{}, err
		}
		users = append(users, u)
	}

	return users, nil
}

// Implements the UserRepository.UpdateUser method
func (pg *Pg) UpdateUser(ctx context.Context, id uuid.UUID, req user.UpdateUserRequest) (user.User, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return user.User{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		UPDATE users 
		SET 
			email = COALESCE($1, email),
			password_hash = COALESCE($2, password_hash),
			updated_at = $3 
		WHERE id = $4
		RETURNING id, email, created_at, updated_at
	`, req.Email, req.PasswordHash, time.Now().UTC(), id)

	var u user.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return user.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Implements the TaskRepository.CreateListTask method
func (pg *Pg) CreateListTask(ctx context.Context, listID uuid.UUID, req task.CreateTaskRequest) (task.Task, error) {
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

	_, err = pg.db.ExecContext(ctx, `
		INSERT INTO lists_tasks (list_id, task_id)
		VALUES ($1, $2)
	`, listID, taskId)
	if err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskRepository.GetListTask method
func (pg *Pg) GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (task.Task, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return task.Task{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		SELECT id, title, completed, created_at, updated_at FROM tasks WHERE id = $1
	`, taskID)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskRepository.GetListTasks method
func (pg *Pg) GetListTasks(ctx context.Context, listID uuid.UUID) ([]task.Task, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, title, completed, created_at, updated_at FROM tasks
		JOIN lists_tasks ON tasks.id = lists_tasks.task_id
		WHERE lists_tasks.list_id = $1
	`, listID)
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

// Implements the TaskRepository.UpdateListTask method
func (pg *Pg) UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req task.UpdateTaskRequest) (task.Task, error) {
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
	`, req.Title, req.Completed, time.Now().UTC(), taskID)

	var t task.Task
	if err := row.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt, &t.Completed); err != nil {
		return task.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the ListRepository.CreateList method
func (pg *Pg) CreateList(ctx context.Context, req list.CreateListRequest) (list.List, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return list.List{}, err
	}
	defer tx.Rollback()

	listId := uuid.New()
	row := pg.db.QueryRowContext(ctx, `
		INSERT INTO lists (id, name)
		VALUES ($1, $2)
		RETURNING id, name, created_at, updated_at
	`, listId, req.Name)

	var l list.List
	if err := row.Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt); err != nil {
		return list.List{}, err
	}

	if err := tx.Commit(); err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListRepository.UpdateList method
func (pg *Pg) UpdateList(ctx context.Context, id uuid.UUID, req list.UpdateListRequest) (list.List, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return list.List{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		UPDATE lists SET name = $1, updated_at = $2 WHERE id = $3
		RETURNING id, name, created_at, updated_at
	`, req.Name, time.Now().UTC(), id)

	var l list.List
	if err := row.Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt); err != nil {
		return list.List{}, err
	}

	if err := tx.Commit(); err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListRepository.GetList method
func (pg *Pg) GetList(ctx context.Context, id uuid.UUID) (list.List, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return list.List{}, err
	}
	defer tx.Rollback()

	row := pg.db.QueryRowContext(ctx, `
		SELECT id, name, created_at, updated_at FROM lists WHERE id = $1
	`, id)

	var l list.List
	if err := row.Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt); err != nil {
		return list.List{}, err
	}

	if err := tx.Commit(); err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListRepository.GetLists method
func (pg *Pg) GetLists(ctx context.Context) ([]list.List, error) {
	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, name, created_at, updated_at FROM lists
	`)
	if err != nil {
		return []list.List{}, err
	}
	defer rows.Close()

	var lists []list.List
	for rows.Next() {
		var l list.List
		if err := rows.Scan(&l.ID, &l.Name, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return []list.List{}, err
		}
		lists = append(lists, l)
	}

	return lists, nil
}
