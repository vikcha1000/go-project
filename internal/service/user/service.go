package user

import (
	"context"
	"errors"
	"mine/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// createUser создает нового Юзера
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*model.User, error) {
	user := model.User{
		Name:             req.Name,
		TelegramUsername: req.TelegramUsername,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Получение Юзера по ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Обновление Юзера
func (s *UserService) UpdateUserByID(ctx context.Context, id uint, req UpdateUserRequest) (*model.User, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.TelegramUsername != nil {
		updates["telegram_username"] = *req.TelegramUsername
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := s.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
