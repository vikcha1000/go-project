package handlers

import (
	"mine/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// createTask создает новую задачу
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	task := new(model.Task)

	// Парсим тело запроса
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Устанавливаем значения по умолчанию
	task.IsDone = false
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Сохраняем в базу данных
	result := h.DB.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.JSON(task)
}

// getTaskByID возвращает задачу по ID
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var task model.Task
	result := h.DB.First(&task, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get task",
		})
	}

	return c.JSON(task)
}

// func GetTasks(c *fiber.Ctx) error {
// 	db := database.GetDB()

// 	tasks, err := model.GetAllTasks(db)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}
// 	return c.JSON(tasks)
// }

// func CreateTask(c *fiber.Ctx) error {
// 	db := database.GetDB()
// 	var task model.TaskRequest
// 	if err := c.BodyParser(&task); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	if err := model.CreateTask(db, &task); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(task)

// }

// func CreateUser(c *fiber.Ctx) error {
// 	db := database.GetDB()
// 	var user model.User
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	if err := model.CreateUser(db, &user); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(user)
// }

// func GetUserByID(c *fiber.Ctx) error {
// 	db := database.GetDB()
// 	id := c.Params("id")
// 	user, err := model.GetUserByID(db, id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}
// 	return c.JSON(user)
// }

// func UpdateUserByID(c *fiber.Ctx) error {

// 	db := database.GetDB()
// 	id := c.Params("id")

// 	// Парсим входные данные
// 	var user model.User
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	// Обновляем пользователя
// 	updatedUser, err := model.UpdateUserByID(db, id, &user)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	if updatedUser == nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "User not found",
// 		})
// 	}

// 	return c.JSON(updatedUser)
// }
