package main

import (
	"log"
	"mine/internal/app"
)

func main() {
	if err:= app.Run(); err != nil{
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}