package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAccessAndRefreshTokens(userID int32, accesExpirationTime int64, refreshExpirationTime int64) (string, string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func Expirations() (int64, int64) {
	accesExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	return accesExpirationTime, refreshExpirationTime
}

func NewTestToken(userID int32) (string, error) {
	accesExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func ValidateSession(tokenString string) (int32, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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