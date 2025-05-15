package app

import (
    "github.com/gofiber/fiber/v2"
    "mine/internal/handlers"
)

func setupRoutes(app *fiber.App) {
    api := app.Group("/api")
    items := api.Group("/items")
    items.Get("/", handlers.GetItems)
    items.Post("/", handlers.CreateItem)
}