package ports

import (
	"context"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/google/uuid"
)

type DatabaseService interface {
	StatusCheck(ctx context.Context) error
}

type AuthService interface {
	Login(ctx context.Context, req domain.CreateSessionRequest) (domain.Session, error)
	Logout(ctx context.Context, req domain.DeleteSessionRequest) error
	LogoutByEmail(ctx context.Context, email string) error
	Refresh(ctx context.Context, req domain.RefreshSessionRequest) (domain.Session, error)
	GetSession(ctx context.Context, req domain.GetSessionRequest) (domain.Session, error)
	RevokeSession(ctx context.Context, req domain.RevokeSessionRequest) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, req domain.CreateSessionRequest) (domain.Session, error)
	RefreshSession(ctx context.Context, req domain.RefreshSessionRequest) (domain.Session, error)
	RevokeSession(ctx context.Context, req domain.RevokeSessionRequest) error
	DeleteSession(ctx context.Context, req domain.DeleteSessionRequest) error
	GetSession(ctx context.Context, req domain.GetSessionRequest) (domain.Session, error)
}

type UserService interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (domain.User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (domain.User, error)
}

type TaskService interface {
	CreateListTask(ctx context.Context, listID uuid.UUID, req domain.CreateTaskRequest) (domain.Task, error)
	GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (domain.Task, error)
	GetListTasks(ctx context.Context, listID uuid.UUID) ([]domain.Task, error)
	UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req domain.UpdateTaskRequest) (domain.Task, error)
}

type TaskRepository interface {
	CreateListTask(ctx context.Context, listID uuid.UUID, req domain.CreateTaskRequest) (domain.Task, error)
	GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (domain.Task, error)
	GetListTasks(ctx context.Context, listID uuid.UUID) ([]domain.Task, error)
	UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req domain.UpdateTaskRequest) (domain.Task, error)
}

type ListService interface {
	CreateList(ctx context.Context, req domain.CreateListRequest) (domain.List, error)
	GetList(ctx context.Context, id uuid.UUID) (domain.List, error)
	GetLists(ctx context.Context) ([]domain.List, error)
	UpdateList(ctx context.Context, id uuid.UUID, req domain.UpdateListRequest) (domain.List, error)
}

type ListRepository interface {
	CreateList(ctx context.Context, req domain.CreateListRequest) (domain.List, error)
	GetList(ctx context.Context, id uuid.UUID) (domain.List, error)
	GetLists(ctx context.Context) ([]domain.List, error)
	UpdateList(ctx context.Context, id uuid.UUID, req domain.UpdateListRequest) (domain.List, error)
}
