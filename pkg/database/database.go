package database

import (
	"fmt"
	"log"
	"mine/internal/model"
	"mine/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.LoadDBConfig()
	// Инициализация подключения к PostgreSQL

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автомиграция
	err = DB.AutoMigrate(&model.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
