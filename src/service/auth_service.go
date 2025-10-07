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
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("user with this email already exists")

type AuthService struct {
	queries *repository.Queries
	ctx     context.Context
	jwtKey  []byte
}

type Session struct {
	UserId           int32
	UserUuid         string
	Token            string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}

func NewAuthService(conn *pgxpool.Pool) *AuthService {
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
			return nil, errors.New("invalid email")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	accesExpirationTime, refreshExpirationTime := Expirations()

	accessTokenString, refreshTokenString, err := CreateAccessAndRefreshTokens(user.ID, user.IsAdmin, accesExpirationTime, refreshExpirationTime)

	if err != nil {
		return nil, err
	}

	return &Session{
		UserId:           user.ID,
		UserUuid:         user.Uuid.String(),
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        s.expirationsTime(accesExpirationTime),
		RefreshExpiresAt: s.expirationsTime(refreshExpirationTime),
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
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			return nil, errors.New("invalid token")
		}

		userID := int32(claims["sub"].(float64))

		var isAdmin bool = false
		if claims["is_admin"] != nil {
			isAdmin = claims["is_admin"].(bool)
		}

		accesExpirationTime, refreshExpirationTime := Expirations()

		newAccessTokenString, newRefreshToken, err := CreateAccessAndRefreshTokens(userID, isAdmin, accesExpirationTime, refreshExpirationTime)

		if err != nil {
			return nil, err
		}
		user_uuid, err := s.queries.GetUserById(s.ctx, userID)
		if err != nil {
			return nil, err
		}

		return &Session{
			UserId:           userID,
			UserUuid:         user_uuid.Uuid.String(),
			Token:            newAccessTokenString,
			RefreshToken:     newRefreshToken,
			ExpiresAt:        s.expirationsTime(accesExpirationTime),
			RefreshExpiresAt: s.expirationsTime(refreshExpirationTime),
		}, nil
	}

	return nil, errors.New("invalid refresh token")
}

func (s *AuthService) Signup(name, email, password string) (*Session, error) {
	_, err := s.queries.GetUserByEmail(s.ctx, email)
	if err == nil {
		return nil, ErrUserExists
	}

	if err != pgx.ErrNoRows {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(s.ctx, repository.CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return nil, err
	}

	accesExpirationTime, refreshExpirationTime := Expirations()

	accessTokenString, refreshTokenString, err := CreateAccessAndRefreshTokens(user.ID, false, accesExpirationTime, refreshExpirationTime)

	if err != nil {
		return nil, err
	}

	return &Session{
		UserId:           user.ID,
		UserUuid:         user.Uuid.String(),
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt:        s.expirationsTime(accesExpirationTime),
		RefreshExpiresAt: s.expirationsTime(refreshExpirationTime),
	}, nil
}

func (s *AuthService) expirationsTime(expiresAt int64) time.Time {
	pgExpiresAt := pgtype.Timestamptz{
		Time:  time.Unix(expiresAt, 0),
		Valid: true,
	}
	return pgExpiresAt.Time
}
