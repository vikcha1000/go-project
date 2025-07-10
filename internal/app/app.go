package app

import (
	"mine/internal/service/login"
	"mine/internal/service/task"
	"mine/internal/service/user"
	"mine/pkg/database"
	"os"

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
	defer logger.Sync()

	// 2. Инициализация БД
	if err := database.InitDB(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
		return err
	}

	// 3. Инициализация сервисов
	taskService := task.NewTaskService(database.GetDB())
	userService := user.NewUserService(database.GetDB())
	jwtSecret := os.Getenv("JWT_SECRET")

	loginService := login.NewLoginService(database.GetDB(), jwtSecret)

	// 4. Создание Fiber приложения
	app := fiber.New()

	// 5. Инициализация обработчиков
	handlers := []FeatureHandler{
		task.NewTaskHandler(taskService, userService, logger, jwtSecret),
		user.NewUserHandler(userService, logger),
		login.NewLoginHandler(loginService, logger),
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
