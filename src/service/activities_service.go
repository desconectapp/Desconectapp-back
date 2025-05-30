package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type ActivitiesService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewActivitiesService(conn *pgx.Conn) *ActivitiesService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &ActivitiesService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *ActivitiesService) ListActivities() ([]repository.ActivityRequest, error) {
	activities, err := s.queries.ListActivityRequests(s.ctx)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivitiesService) CreateActivity(params repository.CreateActivityRequestParams) (repository.ActivityRequest, error) {
	activity, err := s.queries.CreateActivityRequest(s.ctx, params)
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	return activity, nil
}
