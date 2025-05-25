package service

import (
	repository "gin/db/generated"
	"context"
	"github.com/jackc/pgx/v5"
	"fmt"
	"os"
)

type Service struct {
	queries *repository.Queries
	ctx    context.Context
}

func NewService() *Service {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	queries := repository.New(conn)
	ctx := context.Background()

	return &Service{
		queries: queries,
		ctx:    ctx,
	}
}

func (s *Service) ListUsers() ([]repository.User, error) {
	users, err := s.queries.ListUsers(s.ctx)
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

func (s *Service) GetUser(userId int32) (repository.User, error) {
	user, err := s.queries.GetUser(s. ctx, userId)
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

func (s *Service) UpdateUser(userParams repository.UpdateUserParams) (int32, error) {
	id, err := s.queries.UpdateUser(s.ctx, userParams)
	if err != nil {
		return -1, err
	}
	return id, nil
}