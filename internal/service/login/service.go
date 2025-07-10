package login

import (
	"context"
	"errors"
	"mine/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService struct {
	db     *gorm.DB
	secret string
}

func NewLoginService(db *gorm.DB, secret string) *LoginService { // Добавьте secret
	return &LoginService{
		db:     db,
		secret: secret,
	}
}

func (s *LoginService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var user models.User
	var login LoginResponse
	if err := s.db.WithContext(ctx).
		Where("telegram_username = ?", req.TelegramUsername).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("invalid password")
	}

	// Генерируем JWT токен
	token, err := generateJWTToken(user.TelegramUsername, s.secret)
	if err != nil {
		return nil, err
	}

	login = LoginResponse{
		TelegramUsername: req.TelegramUsername,
		Token:            token,
	}
	return &login, nil
}
