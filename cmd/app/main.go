package main

import (
	"log"
	"mine/internal/app"
)

func main() {
	// Инициализация и запуск приложения
	if err := app.Run(); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
