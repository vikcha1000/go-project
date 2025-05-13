package app

import (
	"github.com/gofiber/fiber/v2"
)

func Run() error {
	app := fiber.New()
	SetupRoutes(app)
	return app.Listen(":3000")
}
