package app

import (
	"mine/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskHandler *handlers.TaskHandler) {
	api := app.Group("/api")
	task := api.Group("/task")
	task.Post("/", taskHandler.CreateTask)
	task.Get("/:id", taskHandler.GetTaskByID)

	// user := api.Group("/user")
	// user.Post("/", handlers.CreateUser)
	// user.Get("/:id", handlers.GetUserByID)
	// user.Put("/:id", handlers.UpdateUserByID)
	// //user.Delete("/:id", handlers.DeleteUserByID)
}
