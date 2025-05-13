package controllers

// Добавьте ваши обработчики контроллеров здесьpackage config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	// Добавьте ваши конфиги
}
