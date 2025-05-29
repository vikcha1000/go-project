package task

import (
	"mine/internal/model"

	"context"

	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// createTask создает новую задачу
func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest) (*model.Task, error) {
	task := model.Task{
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

//Получение задачи по id
func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	if err := s.db.WithContext(ctx).Preload("Author").Preload("Executor").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

//Валидация на существование связанного юзера: Автора и Исполнителя
func (s *TaskService) ValidateUsersExist(ctx context.Context, authorID, executorID uint) error {
	if err := s.db.WithContext(ctx).First(&model.User{}, authorID).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).First(&model.User{}, executorID).Error; err != nil {
		return err
	}
	return nil
}