package app

import (
	"mine/internal/service/task" // Пакет с TaskHandler
	"mine/internal/service/user" // Пакет с UserHandler
	"mine/pkg/database"          // Пакет с инициализацией БД

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// FeatureHandler - интерфейс для API
type FeatureHandler interface {
	SetupAPI(r fiber.Router)
}

func Run() error {
	// 1. Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync() // Не забываем закрыть логгер

	// 2. Инициализация БД
	if err := database.InitDB(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
		return err
	}

	// 3. Инициализация сервисов
	taskService := task.NewTaskService(database.GetDB())
    userService := user.NewUserService(database.GetDB())

	// 4. Создание Fiber приложения
	app := fiber.New()

	// // Создаем роутер для API (префикс /api)
	// apiRouter := app.Group("/api")

	// 5. Инициализация обработчиков
	handlers := []FeatureHandler{
		task.NewTaskHandler(taskService, logger),
		user.NewUserHandler(userService, logger),
	}

	// 6. Настройка маршрутов
	apiRouter := app.Group("/api")
	for _, handler := range handlers {
		handler.SetupAPI(apiRouter)
	}

	// 7. Запуск сервера
	logger.Info("Starting server on :3000")
	return app.Listen(":3000")
}
