package model

type User struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Name             string `gorm:"size:255;not null" json:"name"`
	TelegramUsername string `gorm:"size:255;not null" json:"telegram_username"`
}