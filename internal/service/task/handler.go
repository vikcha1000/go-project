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

//CreateTask создает и возвращает задачу
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.service.ValidateUsersExist(c.Context(), req.AuthorID, req.ExecutorID); err != nil {
		h.log.Warn("User validation failed", 
			zap.Uint("authorID", req.AuthorID),
			zap.Uint("executorID", req.ExecutorID))
		return fiber.NewError(fiber.StatusBadRequest, "User not found")
	}

	task, err := h.service.CreateTask(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to create task", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(ToTaskResponse(*task))
}

//GetTaskByID возвращает задачу по id
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}

	task, err := h.service.GetTaskByID(c.Context(), uint(id))
	if err != nil {
		h.log.Warn("Task not found", zap.Uint("id", uint(id)))
		return fiber.NewError(fiber.StatusNotFound, "Task not found")
	}

	return c.JSON(ToTaskResponse(*task))
}