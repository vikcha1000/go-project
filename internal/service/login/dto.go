package login

type LoginRequest struct {
	TelegramUsername string `json:"telegramUsername" validate:"required,max=255"`
	Password         string `json:"password" validate:"required,max=10,min=4"`
}

type LoginResponse struct {
	TelegramUsername string `json:"telegramUsername" validate:"required"`
	Token            string `json:"token"`
}
