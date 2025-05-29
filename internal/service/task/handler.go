package task

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service  *TaskService
	validate *validator.Validate
	log      *zap.Logger
}

func NewTaskHandler(service *TaskService, log *zap.Logger) *TaskHandler {
	return &TaskHandler{
		service:  service,
		validate: validator.New(),
		log:      log,
	}
}

func (h *TaskHandler) SetupAPI(r fiber.Router) {
	group := r.Group("/task")
	group.Post("/", h.CreateTask)
	group.Get("/:id", h.GetTaskByID)
}

// CreateTask создает и возвращает задачу
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Некорректное тело запроса")
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.service.ValidateUsersExist(c.Context(), req.AuthorID); err != nil {
		h.log.Warn("User validation failed: author not exist",
			zap.Uint("authorID", req.AuthorID))
		return fiber.NewError(fiber.StatusBadRequest, "Автор - несуществующий юзер")
	}

	if err := h.service.ValidateUsersExist(c.Context(), req.ExecutorID); err != nil {
		h.log.Warn("User validation failed: executor not exist",
			zap.Uint("executorID", req.ExecutorID))
		return fiber.NewError(fiber.StatusBadRequest, "Исполнитель - несуществующий юзер")
	}

	task, err := h.service.CreateTask(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to create task", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(ToTaskResponse(*task))
}

// GetTaskByID возвращает задачу по id
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.log.Warn("Invalid task ID")
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный id")
	}

	task, err := h.service.GetTaskByID(c.Context(), uint(id))
	if err != nil {
		h.log.Warn("Task not found", zap.Uint("id", uint(id)))
		return fiber.NewError(fiber.StatusNotFound, "Задача не найдена")
	}

	return c.JSON(ToTaskResponse(*task))
}
