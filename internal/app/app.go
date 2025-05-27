package app

import (
	"mine/internal/handlers"
	"mine/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func Run() error {

	// Инициализация БД
	database.InitDB()

	// Инициализация Fiber
	app := fiber.New()

	// Инициализируем обработчики с подключением к БД
	taskHandler := handlers.NewTaskHandler(database.DB)

		// Инициализируем обработчики с подключением к БД
	userHandler := handlers.NewUserHandler(database.DB)

	// Настраиваем маршруты
	SetupRoutes(app, taskHandler, userHandler)

	return app.Listen(":3000")
}
