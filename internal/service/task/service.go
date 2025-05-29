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

func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*model.Task, error) {
	var task model.Task
	if err := s.db.WithContext(ctx).Preload("Author").Preload("Executor").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) ValidateUsersExist(ctx context.Context, authorID, executorID uint) error {
	if err := s.db.WithContext(ctx).First(&model.User{}, authorID).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).First(&model.User{}, executorID).Error; err != nil {
		return err
	}
	return nil
}

// func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
// 	var req CreateTaskRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	// Валидация
// 	if err := h.Validate.Struct(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	if err := h.DB.First(&model.User{}, req.AuthorID).Error; err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Author not found",
// 		})
// 	}

// 	task := model.Task{
// 		Name:        req.Name,
// 		Description: req.Description,
// 		AuthorID:    req.AuthorID,
// 		ExecutorID:  req.ExecutorID,
// 		Deadline:    req.Deadline,
// 		IsDone:      false,
// 	}

// 	if err := h.DB.Create(&task).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to create task",
// 		})
// 	}

// 	// Загружаем связанные данные
// 	h.DB.Preload("Author").Preload("Executor").First(&task, task.ID)

// 	return c.Status(fiber.StatusCreated).JSON(ToTaskResponse(task))
// }
