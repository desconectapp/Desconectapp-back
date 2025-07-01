package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type ActivitiesRequestService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewActivitiesRequestService(conn *pgx.Conn) *ActivitiesRequestService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &ActivitiesRequestService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *ActivitiesRequestService) ListActivitiesRequests(params repository.ListActivityRequestsParams) ([]repository.ActivityRequest, error) {
	activities, err := s.queries.ListActivityRequests(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivitiesRequestService) CreateActivityRequest(params repository.CreateActivityRequestParams) (repository.ActivityRequest, error) {
	activity, err := s.queries.CreateActivityRequest(s.ctx, params)
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	return activity, nil
}

func (s *ActivitiesRequestService) GetActivities(params repository.GetActivitiesParams) ([]repository.Activity, error) {
	activities, err := s.queries.GetActivities(s.ctx, params)

	if err != nil {
		return nil, err
	}
	return activities, nil
}