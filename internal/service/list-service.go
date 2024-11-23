package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/list"
	"github.com/charlieroth/reminders/internal/ports"
	"github.com/google/uuid"
)

type ListService struct {
	repo ports.ListRepository
}

func NewListService(repo ports.ListRepository) *ListService {
	return &ListService{repo: repo}
}

// Implements the ListService.CreateList method
func (s *ListService) CreateList(ctx context.Context, req list.CreateListRequest) (list.List, error) {
	l, err := s.repo.CreateList(ctx, req)
	if err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListService.UpdateList method
func (s *ListService) UpdateList(ctx context.Context, id uuid.UUID, req list.UpdateListRequest) (list.List, error) {
	l, err := s.repo.UpdateList(ctx, id, req)
	if err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListService.GetList method
func (s *ListService) GetList(ctx context.Context, id uuid.UUID) (list.List, error) {
	l, err := s.repo.GetList(ctx, id)
	if err != nil {
		return list.List{}, err
	}

	return l, nil
}

// Implements the ListService.GetLists method
func (s *ListService) GetLists(ctx context.Context) ([]list.List, error) {
	lists, err := s.repo.GetLists(ctx)
	if err != nil {
		return []list.List{}, err
	}

	return lists, nil
}
