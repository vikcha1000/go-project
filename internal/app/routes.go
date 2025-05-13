package app

import (
	"mine/internal/handlers" // Исправлено

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	data := api.Group("/data")
	data.Get("/", handlers.GetData)
	data.Post("/", handlers.CreateData)
}
