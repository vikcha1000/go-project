package user

import (
	"mine/internal/model"
	"context"
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
		Name:        req.Name,
		TelegramUsername: req.TelegramUsername,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
    if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
	return &user, nil
}
