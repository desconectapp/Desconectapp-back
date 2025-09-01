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
	UserId 				 int32
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
			return nil, errors.New("invalid email")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	accesExpirationTime, refreshExpirationTime := s.expirations()

	accessTokenString, refreshTokenString, err := s.createAccessAndRefreshTokens(user.ID, accesExpirationTime, refreshExpirationTime)

	if err != nil {
		return nil, err
	}

	return &Session{
		UserId:				user.ID,
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt: 		s.expirationsTime(accesExpirationTime),
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

		accesExpirationTime, refreshExpirationTime := s.expirations()

		newAccessTokenString, newRefreshToken, err := s.createAccessAndRefreshTokens(userID, accesExpirationTime, refreshExpirationTime)

		if err != nil {
			return nil, err
		}

		return &Session{
			UserId:			  userID,
			Token:            newAccessTokenString,
			RefreshToken:     newRefreshToken,
			ExpiresAt: 		s.expirationsTime(accesExpirationTime),
		RefreshExpiresAt: s.expirationsTime(refreshExpirationTime),
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
		exp := int64(claims["exp"].(float64))
			if time.Now().Unix() > exp {
				return 0, errors.New("invalid token")
			}

		return int32(claims["sub"].(float64)), nil
	}

	return 0, errors.New("invalid token")
}

func (s *AuthService) Signup(name, email, password string) (*Session, error) {
	_, err := s.queries.GetUserByEmail(s.ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
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

	accesExpirationTime, refreshExpirationTime := s.expirations()

	accessTokenString, refreshTokenString, err := s.createAccessAndRefreshTokens(user.ID, accesExpirationTime, refreshExpirationTime)

	if err != nil {
		return nil, err
	}

	return &Session{
		UserId:				user.ID,
		Token:            accessTokenString,
		RefreshToken:     refreshTokenString,
		ExpiresAt: 		s.expirationsTime(accesExpirationTime),
		RefreshExpiresAt: s.expirationsTime(refreshExpirationTime),
	}, nil
}

func (s *AuthService) createAccessAndRefreshTokens(userID int32, accesExpirationTime int64, refreshExpirationTime int64) (string, string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(s.jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) expirations() (int64, int64) {
	accesExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	return accesExpirationTime, refreshExpirationTime
}

func (s *AuthService) expirationsTime(expiresAt int64) time.Time {
	pgExpiresAt := pgtype.Timestamptz{
		Time:  time.Unix(expiresAt, 0),
		Valid: true,
	}
	return pgExpiresAt.Time
}
