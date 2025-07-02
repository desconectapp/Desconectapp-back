package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type PreferenceService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewPreferenceService(conn *pgx.Conn) *PreferenceService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &PreferenceService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *PreferenceService) GetUserPreferences(params repository.GetUserPreferencesParams) ([]repository.GetUserPreferencesRow, error ) {
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

func (s *PreferenceService) AddPreference(params repository.AddPreferenceParams) (error) {
	err := s.queries.AddPreference(s.ctx, params) 
	return err
}

func (s *PreferenceService) DeletePreference(params repository.DeletePreferenceParams) (error) {
	err := s.queries.DeletePreference(s.ctx, params)

	return err
}