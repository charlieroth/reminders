package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/ports"
	"github.com/charlieroth/reminders/internal/task"
	"github.com/google/uuid"
)

type Service struct {
	repo ports.TaskRepository
}

func NewService(repo ports.TaskRepository) *Service {
	return &Service{repo: repo}
}

// Implements the TaskService.CreateTask method
func (s *Service) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	t, err := s.repo.CreateTask(ctx, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.ListTasks method
func (s *Service) ListTasks(ctx context.Context) ([]task.Task, error) {
	tasks, err := s.repo.ListTasks(ctx)
	if err != nil {
		return []task.Task{}, err
	}

	return tasks, nil
}

// Implements the TaskService.GetTask method
func (s *Service) GetTask(ctx context.Context, id uuid.UUID) (task.Task, error) {
	t, err := s.repo.GetTask(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.UpdateTask method
func (s *Service) UpdateTask(ctx context.Context, id uuid.UUID, req task.UpdateTaskRequest) (task.Task, error) {
	t, err := s.repo.UpdateTask(ctx, id, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}
