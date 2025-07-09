package login

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWTToken(telegramUsername, secret string) (string, error) {
	claims := jwt.MapClaims{
		//	"user_id":  userID,
		"telegramUsername": telegramUsername,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
