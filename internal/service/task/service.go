package task

import (
	"context"
	"mine/internal/models"

	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// createTask создает новую задачу
func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*models.Task, error) {
	task := models.Task{
		Name:        req.Name,
		Description: req.Description,
		AuthorID:    req.AuthorID,
		ExecutorID:  req.ExecutorID,
		Deadline:    req.Deadline,
		IsDone:      false,
	}

	if err := s.db.WithContext(ctx).Create(&task).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Preload("Author").Preload("Executor").First(&task, task.ID).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// Получение задачи по id
func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	if err := s.db.WithContext(ctx).Preload("Author").Preload("Executor").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Валидация на существование связанного юзера: Автора и Исполнителя
func (s *TaskService) ValidateUsersExist(ctx context.Context, userID uint) error {
	if err := s.db.WithContext(ctx).First(&models.User{}, userID).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).First(&models.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

// Обновление задачи
func (s *TaskService) UpdateTaskByID(ctx context.Context, id uint, req UpdateTaskRequest) (*models.Task, error) {
	var task models.Task
	return &task, nil
}
