package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/ports"
	"github.com/charlieroth/reminders/internal/task"
)

type Service struct {
	repo ports.TaskRepository
}

func NewService(repo ports.TaskRepository) *Service {
	return &Service{repo: repo}
}

// Implement the TaskService interface
func (s *Service) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.Task, error) {
	t, err := s.repo.CreateTask(ctx, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}
