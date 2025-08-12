package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewService(conn *pgx.Conn) *Service {
	queries := repository.New(conn)
	ctx := context.Background()

	return &Service{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *Service) ListUsers(params repository.ListUsersParams) ([]repository.User, error) {
	users, err := s.queries.ListUsers(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) CreateUser(userParams repository.CreateUserParams) (repository.User, error) {
	user, err := s.queries.CreateUser(s.ctx, userParams)
	if err != nil {
		return repository.User{}, err
	}
	return user, nil
}

func (s *Service) CreateProfile(profile repository.CreateProfileParams) (int32, error) {
	id, err := s.queries.CreateProfile(s.ctx, profile)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) GetUser(userId int32) (repository.User, error) {
	user, err := s.queries.GetUser(s.ctx, userId)
	if err != nil {
		return repository.User{}, err
	}
	return user, nil
}

func (s *Service) DeleteUser(userId int32) (int32, error) {
	id, err := s.queries.DeleteUser(s.ctx, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) UpdateUser(userParams repository.UpdateUserParams, userPreferences repository.BatchAddPreferencesParams) (int32, error) {
	id, err := s.queries.UpdateUser(s.ctx, userParams)
	if err != nil {
		return -1, err
	}
	err = s.queries.BatchAddPreferences(s.ctx, userPreferences)
	if err != nil {
		return -1, err
	}
	return id, nil
}