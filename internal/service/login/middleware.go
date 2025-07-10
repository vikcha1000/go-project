package login

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Token:", tokenString) // Логируем токен

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			fmt.Println("Token parse error:", err) // Логируем ошибки парсинга
			return c.Status(401).JSON("тест")
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
