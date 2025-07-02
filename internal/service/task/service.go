package task

import (
	"context"
	"errors"
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
	return nil
}

// Обновление задачи
func (s *TaskService) UpdateTaskByID(ctx context.Context, id uint, req UpdateTaskRequest) (*models.Task, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.AuthorID != nil {
		updates["author_id"] = *req.AuthorID
	}
	if req.ExecutorID != nil {
		updates["executor_id"] = *req.ExecutorID
	}
	if req.IsDone != nil {
		updates["is_done"] = *req.IsDone
	}
	if req.Deadline != nil {
		updates["deadline"] = *req.Deadline
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}
	result := s.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var task models.Task
	if err := s.db.WithContext(ctx).Preload("Author").Preload("Executor").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Получение задач по ExecutorID
func (s *TaskService) GetTasksByExecutorID(ctx context.Context, executorID uint) ([]*models.Task, error) {
	var tasks []*models.Task

	err := s.db.WithContext(ctx).
		Preload("Author").
		Preload("Executor").
		Where("executor_id = ?", executorID).
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}
