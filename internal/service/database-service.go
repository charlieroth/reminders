package service

import (
	"context"

	"github.com/charlieroth/reminders/internal/outbound"
)

type DatabaseService struct {
	pg *outbound.Pg
}

func NewDatabaseService(pg *outbound.Pg) *DatabaseService {
	return &DatabaseService{pg: pg}
}

func (ds *DatabaseService) StatusCheck(ctx context.Context) error {
	return ds.pg.StatusCheck(ctx)
}
