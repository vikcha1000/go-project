package user

// --- Response DTOs ---

type UserResponse struct {
	ID               uint   `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	TelegramUsername string `json:"telegramUsername" validate:"required"`
}

type CreateUserRequest struct {
	Name             string `json:"name" validate:"required,min=3"`
	TelegramUsername string `json:"telegramUsername" validate:"required,max=255"`
}

type UpdateUserRequest struct {
	Name             *string `json:"name,omitempty" validate:"omitempty,min=3"`
	TelegramUsername *string `json:"telegramUsername,omitempty" validate:"omitempty,max=255"`
}
