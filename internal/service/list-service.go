package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/domain"
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
func (s *ListService) CreateList(ctx context.Context, req domain.CreateListRequest) (domain.List, error) {
	l, err := s.repo.CreateList(ctx, req)
	if err != nil {
		return domain.List{}, err
	}

	return l, nil
}

// Implements the ListService.UpdateList method
func (s *ListService) UpdateList(ctx context.Context, id uuid.UUID, req domain.UpdateListRequest) (domain.List, error) {
	l, err := s.repo.UpdateList(ctx, id, req)
	if err != nil {
		return domain.List{}, err
	}

	return l, nil
}

// Implements the ListService.GetList method
func (s *ListService) GetList(ctx context.Context, id uuid.UUID) (domain.List, error) {
	l, err := s.repo.GetList(ctx, id)
	if err != nil {
		return domain.List{}, err
	}

	return l, nil
}

// Implements the ListService.GetLists method
func (s *ListService) GetLists(ctx context.Context) ([]domain.List, error) {
	lists, err := s.repo.GetLists(ctx)
	if err != nil {
		return []domain.List{}, err
	}

	return lists, nil
}
