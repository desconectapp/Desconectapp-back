package service

import (
	"context"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5"
)

type ChatsService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewChatsService(conn *pgx.Conn) *ChatsService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &ChatsService{
		queries: queries,
		ctx:     ctx,
	}
}
