package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secret)
}

func GenerateAccess(userID int, secret string, ttl time.Duration) (string, error) {
	return generate(userID, secret, ttl)
}

func GenerateRefresh(userID int, secret string, ttl time.Duration) (string, error) {
	return generate(userID, secret, ttl)
}

func generate(userID int, secret string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func Parse(tokenString string, secret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
