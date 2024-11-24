package ports

import (
	"context"

	"github.com/charlieroth/reminders/internal/list"
	"github.com/charlieroth/reminders/internal/task"
	"github.com/charlieroth/reminders/internal/user"
	"github.com/google/uuid"
)

type DatabaseService interface {
	StatusCheck(ctx context.Context) error
}

type UserService interface {
	CreateUser(ctx context.Context, req user.CreateUserRequest) (user.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
	GetUsers(ctx context.Context) ([]user.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req user.UpdateUserRequest) (user.User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, req user.CreateUserRequest) (user.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
	GetUsers(ctx context.Context) ([]user.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req user.UpdateUserRequest) (user.User, error)
}

type TaskService interface {
	CreateListTask(ctx context.Context, listID uuid.UUID, req task.CreateTaskRequest) (task.Task, error)
	GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (task.Task, error)
	GetListTasks(ctx context.Context, listID uuid.UUID) ([]task.Task, error)
	UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req task.UpdateTaskRequest) (task.Task, error)
}

type TaskRepository interface {
	CreateListTask(ctx context.Context, listID uuid.UUID, req task.CreateTaskRequest) (task.Task, error)
	GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (task.Task, error)
	GetListTasks(ctx context.Context, listID uuid.UUID) ([]task.Task, error)
	UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req task.UpdateTaskRequest) (task.Task, error)
}

type ListService interface {
	CreateList(ctx context.Context, req list.CreateListRequest) (list.List, error)
	GetList(ctx context.Context, id uuid.UUID) (list.List, error)
	GetLists(ctx context.Context) ([]list.List, error)
	UpdateList(ctx context.Context, id uuid.UUID, req list.UpdateListRequest) (list.List, error)
}

type ListRepository interface {
	CreateList(ctx context.Context, req list.CreateListRequest) (list.List, error)
	GetList(ctx context.Context, id uuid.UUID) (list.List, error)
	GetLists(ctx context.Context) ([]list.List, error)
	UpdateList(ctx context.Context, id uuid.UUID, req list.UpdateListRequest) (list.List, error)
}
