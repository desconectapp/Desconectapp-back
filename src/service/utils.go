package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAccessAndRefreshTokens(userID int32, is_admin bool, accesExpirationTime int64, refreshExpirationTime int64) (string, string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	body := jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	}
	if is_admin {
		body["is_admin"] = true
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, body)

	body = jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpirationTime,
	}
	if is_admin {
		body["is_admin"] = true
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, body)

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

func ValidateSession(tokenString string) (int32, bool, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return -1, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			return -2, false, errors.New("invalid token")
		}

		isAdmin := claims["is_admin"].(bool)

		return int32(claims["sub"].(float64)), isAdmin, nil
	}

	return -3, false, errors.New("invalid token")
}

