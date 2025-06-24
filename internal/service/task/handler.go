package task

import (
	"errors"
	"mine/pkg/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	group.Put("/:id", h.UpdateTaskByID)
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

	if err := h.service.ValidateUsersExist(c.Context(), req.AuthorID); err != nil {
		return errs.Error(c, errs.ErrAuthorNotExist, nil)
	}

	if err := h.service.ValidateUsersExist(c.Context(), req.ExecutorID); err != nil {
		return errs.Error(c, errs.ErrexecutorNotExist, nil)
	}

	task, err := h.service.CreateTask(c.Context(), req)
	if err != nil {
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, task, "")
}

// GetTaskByID возвращает задачу по id
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	task, err := h.service.GetTaskByID(c.Context(), uint(id))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.Error(c, errs.ErrNotFound, nil)
	}

	if err != nil {
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, task, "")
}

// UpdateTaskByID обновляет и возвращает задачу по id
func (h *TaskHandler) UpdateTaskByID(c *fiber.Ctx) error {
	var req UpdateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	task, err := h.service.UpdateTaskByID(c.Context(), uint(id), req)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.Error(c, errs.ErrNotFound, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.service.ValidateUsersExist(c.Context(), req.AuthorID); err != nil {
		return errs.Error(c, errs.ErrAuthorNotExist, nil)
	}

	// 	if err := h.service.ValidateUsersExist(c.Context(), req.ExecutorID); err != nil {
	// 		return errs.Error(c, errs.ErrexecutorNotExist, nil)
	// 	}

	return errs.Success(c, task, "")
}
