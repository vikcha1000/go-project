package task

import (
	"errors"
	"mine/internal/service/login"
	"mine/internal/service/user"
	"mine/pkg/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskHandler struct {
	taskService *TaskService
	userService *user.UserService
	validate    *validator.Validate
	log         *zap.Logger
	jwtSecret   string
}

func NewTaskHandler(taskService *TaskService, userService *user.UserService, log *zap.Logger, jwtSecret string) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		userService: userService,
		validate:    validator.New(),
		log:         log,
		jwtSecret:   jwtSecret,
	}
}

func (h *TaskHandler) SetupAPI(r fiber.Router) {

	authMiddleware := login.AuthMiddleware(h.jwtSecret)

	groupTasks := r.Group("/tasks")
	groupTasks.Use(authMiddleware) // Применяем middleware
	groupTasks.Get("/author", h.GetAuthorTasks)

	groupTask := r.Group("/task")
	groupTask.Post("/", h.CreateTask)
	groupTask.Get("/:id", h.GetTaskByID)
	groupTask.Put("/:id", h.UpdateTaskByID)
	groupTasks.Get("/", h.GetTasksByExecutorId)
	groupTasks.Get("/author", h.GetAuthorTasks)
}

// CreateTask создает и возвращает задачу
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.taskService.ValidateUsersExist(c.Context(), req.AuthorID); err != nil {
		return errs.Error(c, errs.ErrAuthorNotExist, nil)
	}

	if err := h.taskService.ValidateUsersExist(c.Context(), req.ExecutorID); err != nil {
		return errs.Error(c, errs.ErrExecutorNotExist, nil)
	}

	task, err := h.taskService.CreateTask(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to create task", zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, task, "")
}

// GetTaskByID возвращает задачу по id
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	task, err := h.taskService.GetTaskByID(c.Context(), uint(id))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.Error(c, errs.ErrNotFound, nil)
	}

	if err != nil {
		h.log.Error("Failed to get task", zap.Uint("id", uint(id)), zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, task, "")
}

// UpdateTaskByID обновляет и возвращает задачу по id
func (h *TaskHandler) UpdateTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}
	var req UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if req.Name == nil && req.Description == nil && req.AuthorID == nil && req.ExecutorID == nil && req.IsDone == nil && req.Deadline == nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}
	if req.AuthorID != nil {
		if err := h.taskService.ValidateUsersExist(c.Context(), *req.AuthorID); err != nil {
			return errs.Error(c, errs.ErrAuthorNotExist, nil)
		}
	}
	if req.ExecutorID != nil {
		if err := h.taskService.ValidateUsersExist(c.Context(), *req.ExecutorID); err != nil {
			return errs.Error(c, errs.ErrExecutorNotExist, nil)
		}
	}
	task, err := h.taskService.UpdateTaskByID(c.Context(), uint(id), req)
	if err != nil {
		switch {

		case errors.Is(err, gorm.ErrRecordNotFound):
			return errs.Error(c, errs.ErrNotFound, nil)

		case errors.Is(err, errors.New("no fields to update")):
			return errs.Error(c, errs.ErrEmptyBody, nil)
		default:
			h.log.Error("Failed to update task", zap.Uint("id", uint(id)), zap.Error(err))
			return errs.Error(c, errs.ErrInternal, nil)
		}
	}

	return errs.Success(c, task, "")
}

// GetTasksByExecutorId возвращает задачи по ExecutorId
func (h *TaskHandler) GetTasksByExecutorId(c *fiber.Ctx) error {
	executorIdStr := c.Query("executorId")
	if executorIdStr == "" {
		return errs.Error(c, errs.ErrExecutorIdEmpty, nil)
	}
	executorId := c.QueryInt("executorId")
	if executorId == 0 || executorId <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}
	tasks, err := h.taskService.GetTasksByExecutorID(c.Context(), uint(executorId))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.Error(c, errs.ErrNotFound, nil)
	}

	if err := h.taskService.ValidateUsersExist(c.Context(), uint(executorId)); err != nil {
		return errs.Error(c, errs.ErrExecutorNotExist, nil)
	}

	if err != nil {
		h.log.Error("Failed to get tasks by executorId", zap.Uint("executorId", uint(executorId)), zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, tasks, "")

}

// GetAuthorTasks возвращает задачи по AuthorId авторизованного Юзера
func (h *TaskHandler) GetAuthorTasks(c *fiber.Ctx) error {
	// Получаем username из токена
	telegramUsername, ok := c.Locals("telegramUsername").(string)
	if !ok {
		h.log.Error("Telegram username not found in context")
		h.log.Error(telegramUsername)
		return errs.Error(c, errs.ErrTelegramUernameInToken, nil)
	}

	// Получаем пользователя
	user, err := h.userService.GetUserByTelegramUserName(c.Context(), telegramUsername)
	if err != nil {
		h.log.Error("Failed to get user", zap.Error(err))
		return errs.Error(c, errs.ErrTelegramUernameNotExists, nil)
	}

	// Получаем задачи
	tasks, err := h.taskService.GetAuthorTasks(c.Context(), user.ID)
	if err != nil {
		h.log.Error("Failed to get tasks", zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, tasks, "")
}
