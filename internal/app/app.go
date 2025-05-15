package app

import (
    "log"
    "mine/internal/model"
    "mine/pkg/database"
    "github.com/gofiber/fiber/v2"
)

func Run() error {
    // Инициализация БД
    if err := database.InitDB(); err != nil {
        log.Fatalf("Database initialization failed: %v", err)
    }

    // Создаем таблицу
    db := database.GetDB()
    if err := models.CreateTable(db); err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }

    app := fiber.New()
    setupRoutes(app)
    
    return app.Listen(":3000")
}