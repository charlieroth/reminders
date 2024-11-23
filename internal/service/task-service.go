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

// Implements the TaskService.CreateListTask method
func (s *TaskService) CreateListTask(ctx context.Context, listID uuid.UUID, req task.CreateTaskRequest) (task.Task, error) {
	t, err := s.repo.CreateListTask(ctx, listID, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.GetListTasks method
func (s *TaskService) GetListTasks(ctx context.Context, listID uuid.UUID) ([]task.Task, error) {
	tasks, err := s.repo.GetListTasks(ctx, listID)
	if err != nil {
		return []task.Task{}, err
	}

	return tasks, nil
}

// Implements the TaskService.GetListTask method
func (s *TaskService) GetListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID) (task.Task, error) {
	t, err := s.repo.GetListTask(ctx, listID, taskID)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

// Implements the TaskService.UpdateListTask method
func (s *TaskService) UpdateListTask(ctx context.Context, listID uuid.UUID, taskID uuid.UUID, req task.UpdateTaskRequest) (task.Task, error) {
	t, err := s.repo.UpdateListTask(ctx, listID, taskID, req)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}
