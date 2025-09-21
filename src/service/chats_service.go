package service

import (
	"context"
	"fmt"
	repository "gin/db/generated"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatsService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewChatsService(conn *pgxpool.Pool) *ChatsService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &ChatsService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *ChatsService) GetToken(userID int32) (string, error) {
	secret := os.Getenv("SUPABASE_JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("SUPABASE_JWT_SECRET not configured")
	}

	user, err := s.queries.GetUserById(s.ctx, userID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":  user.Uuid.String(),                      // auth.uid() will return this
		"exp":  time.Now().Add(15 * time.Minute).Unix(), // short-lived for safety
		"role": "authenticated",                         // optional but common for RLS policies
	}

	// Sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
