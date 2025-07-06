package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type GroupsService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewGroupsService(conn *pgx.Conn) *GroupsService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &GroupsService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *GroupsService) ListGroups(params repository.ListGroupsParams) ([]repository.Group, error) {
	groups, err := s.queries.ListGroups(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

