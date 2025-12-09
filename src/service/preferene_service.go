package service

import (
	"context"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PreferenceService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewPreferenceService(conn *pgxpool.Pool) *PreferenceService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &PreferenceService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *PreferenceService) GetUserPreferences(params repository.GetUserPreferencesParams) ([]repository.GetUserPreferencesRow, error) {
	preferences, err := s.queries.GetUserPreferences(s.ctx, params)

	if err != nil {
		return nil, err
	}
	return preferences, nil
}

func (s *PreferenceService) BatchAddPreferences(params repository.BatchAddPreferencesParams) error {
	err := s.queries.BatchAddPreferences(s.ctx, params)
	return err
}

func (s *PreferenceService) AddPreference(params repository.AddPreferenceParams) error {
	err := s.queries.AddPreference(s.ctx, params)
	return err
}

func (s *PreferenceService) DeletePreference(params repository.DeletePreferenceParams) (int32, error) {
	id, err := s.queries.DeletePreference(s.ctx, params)
	return id, err
}

func (s *PreferenceService) GenerateCustomActivity(name string) (int32, error) {
	activity, err := GenerateActivity(&ActivitiesRequestService{
		queries: s.queries,
		ctx:     s.ctx,
	}, name)
	if err != nil {
		return 0, err
	}
	return activity.ID, nil
}
