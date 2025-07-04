package service

import (
	"context"
	"errors"
	repository "gin/db/generated"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *repository.Queries
	ctx     context.Context
	jwtKey  []byte
}

type Session struct {
	Token            string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}

func NewAuthService(conn *pgx.Conn) *AuthService {
	queries := repository.New(conn)
	ctx := context.Background()
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		panic("JWT_SECRET environment variable is not set")
	}
	return &AuthService{
		queries: queries,
		ctx:     ctx,
		jwtKey:  jwtKey,
	}
}

func (s *AuthService) Login(email, password string) (*Session, error) {
	user, err := s.queries.GetUserByEmail(s.ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	pgExpiresAt := pgtype.Timestamptz{
		Time:  expiresAt,
		Valid: true,
	}
	pgRefreshExpiresAt := pgtype.Timestamptz{
		Time:  refreshExpiresAt,
		Valid: true,
	}

	_, err = s.queries.CreateSession(s.ctx, repository.CreateSessionParams{
		UserID:           user.ID,
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        pgExpiresAt,
		RefreshExpiresAt: pgRefreshExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &Session{
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*Session, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		session, err := s.queries.GetSessionByRefreshToken(s.ctx, refreshToken)
		if err != nil {
			return nil, errors.New("invalid refresh token")
		}

		if !session.RefreshExpiresAt.Valid || session.RefreshExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("refresh token expired")
		}

		userID := int32(claims["sub"].(float64))

		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(1 * time.Hour).Unix(),
		})

		newAccessTokenString, err := newAccessToken.SignedString(s.jwtKey)
		if err != nil {
			return nil, err
		}

		expiresAt := time.Now().Add(1 * time.Hour)
		pgExpiresAt := pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		}

		_, err = s.queries.UpdateSessionToken(s.ctx, repository.UpdateSessionTokenParams{
			Token:     newAccessTokenString,
			ExpiresAt: pgExpiresAt,
			ID:        session.ID,
		})
		if err != nil {
			return nil, err
		}

		return &Session{
			Token:            newAccessTokenString,
			RefreshToken:     refreshToken,
			ExpiresAt:        expiresAt,
			RefreshExpiresAt: session.RefreshExpiresAt.Time,
		}, nil
	}

	return nil, errors.New("invalid refresh token")
}

func (s *AuthService) ValidateSession(tokenString string) (int32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		session, err := s.queries.GetSessionByToken(s.ctx, tokenString)
		if err != nil {
			return 0, err
		}

		if !session.ExpiresAt.Valid || session.ExpiresAt.Time.Before(time.Now()) {
			return 0, errors.New("session expired")
		}

		return int32(claims["sub"].(float64)), nil
	}

	return 0, errors.New("invalid token")
}

func (s *AuthService) Logout(tokenString string) error {
	return s.queries.DeleteSessionByToken(s.ctx, tokenString)
}

func (s *AuthService) Signup(name, email, password string) (*Session, error) {
	_, err := s.queries.GetUserByEmail(s.ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(s.ctx, repository.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	pgExpiresAt := pgtype.Timestamptz{
		Time:  expiresAt,
		Valid: true,
	}
	pgRefreshExpiresAt := pgtype.Timestamptz{
		Time:  refreshExpiresAt,
		Valid: true,
	}

	// Create session
	_, err = s.queries.CreateSession(s.ctx, repository.CreateSessionParams{
		UserID:           user.ID,
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        pgExpiresAt,
		RefreshExpiresAt: pgRefreshExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &Session{
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
