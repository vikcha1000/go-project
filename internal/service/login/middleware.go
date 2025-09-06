package login

import (
	"fmt"
	"mine/pkg/errs"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return errs.Error(c, errs.ErrUnauthorized, nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Token:", tokenString) // Логируем токен

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if strings.Contains(err.Error(), "token is expired") {
			return errs.Error(c, errs.ErrTokenExpired, nil)
		}

		if !token.Valid {
			return errs.Error(c, errs.ErrUnauthorized, nil)
		}
		claims := token.Claims.(jwt.MapClaims)

		if username, ok := claims["telegramUsername"].(string); ok {
			fmt.Println("Setting username:", username) // Логируем username
			c.Locals("telegramUsername", username)
		} else {
			fmt.Println("telegramUsername not found in claims") // Логируем проблему
		}

		return c.Next()
	}
}
