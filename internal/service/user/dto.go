package user

import (

)

// --- Response DTOs ---

type UserResponse struct {
	ID               uint   `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	TelegramUsername string `json:"telegram_username" validate:"required"`
}