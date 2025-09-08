package user

import (
	"context"
	"errors"
	"fmt"
	"mine/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// createUser создает нового Юзера
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	// Проверяем существование пользователя
	var existingUser models.User
	if err := s.db.WithContext(ctx).
		Where("telegram_username = ?", req.TelegramUsername).
		First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("user with telegram username '%s' already exists", req.TelegramUsername)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаём пользователя
	user := models.User{
		Name:             req.Name,
		TelegramUsername: req.TelegramUsername,
		PasswordHash:     string(hashedPassword),
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Получение Юзера по ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// // Получение Авторизованного Юзера
// func (s *UserService) GetMyUser(ctx context.Context, id uint) (*models.User, error) {
// 	var user models.User
// 	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// Обновление Юзера
func (s *UserService) UpdateUserByID(ctx context.Context, id uint, req UpdateUserRequest) (*models.User, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.TelegramUsername != nil {
		var existingUser models.User
		if err := s.db.WithContext(ctx).
			Where("telegram_username = ?", req.TelegramUsername).
			First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("user with telegram username '%s' already exists", *req.TelegramUsername)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		updates["telegram_username"] = *req.TelegramUsername
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password_hash"] = hashedPassword
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Удаляет Юзера, если он не создал задачи или не назначен исполнителем
func (s *UserService) DeleteUserByID(ctx context.Context, id uint) error {
	// Проверяем наличие связанных задач
	var taskCountAuthor, taskCountExecutor int64
	if err := s.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("author_id = ?", id).
		Count(&taskCountAuthor).Error; err != nil {
		return fmt.Errorf("failed to check tasks: %w", err)
	}

	if taskCountAuthor > 0 {
		return errors.New("User has associated tasks: Author")
	}

	if err := s.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("executor_id = ?", id).
		Count(&taskCountExecutor).Error; err != nil {
		return fmt.Errorf("failed to check tasks: %w", err)
	}

	if taskCountExecutor > 0 {
		return errors.New("User has associated tasks: Executor")
	}

	// Удаляем пользователя
	result := s.db.WithContext(ctx).
		Delete(&models.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

// Получение Юзера по TelegramUserName
func (s *UserService) GetUserByTelegramUserName(ctx context.Context, telegramUserName string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Where("telegram_username = ?", telegramUserName).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
