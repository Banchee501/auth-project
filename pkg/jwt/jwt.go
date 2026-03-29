package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("super-secret-key")

func Generate(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secret)
}

func GenerateAccess(userID int) (string, error) {
	return generate(userID, time.Minute*15)
}

func GenerateRefresh(userID int) (string, error) {
	return generate(userID, time.Hour*24*7)
}

func generate(userID int, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	})

	return token.SignedString(secret)
}

func Parse(tokenString string) (int, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, err
	}

	userID := int(claims["user_id"].(float64))

	return userID, nil
}
