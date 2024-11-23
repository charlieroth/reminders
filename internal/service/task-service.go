package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/ports"
	"github.com/charlieroth/reminders/internal/task"
	"github.com/google/uuid"
)

type TaskService struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// Implements the TaskService.CreateTask method
func (s *TaskService) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	t, err := s.repo.CreateTask(ctx, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.ListTasks method
func (s *TaskService) ListTasks(ctx context.Context) ([]task.Task, error) {
	tasks, err := s.repo.ListTasks(ctx)
	if err != nil {
		return []task.Task{}, err
	}

	return tasks, nil
}

// Implements the TaskService.GetTask method
func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (task.Task, error) {
	t, err := s.repo.GetTask(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.UpdateTask method
func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, req task.UpdateTaskRequest) (task.Task, error) {
	t, err := s.repo.UpdateTask(ctx, id, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}
