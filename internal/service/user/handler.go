package user

import (
	"errors"
	"mine/internal/models"
	"mine/pkg/errs"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserHandler struct {
	service  *UserService
	validate *validator.Validate
	log      *zap.Logger
}

func NewUserHandler(service *UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
		log:      log,
	}
}

func (h *UserHandler) SetupAPI(r fiber.Router) {
	user := r.Group("/user")
	user.Post("/", h.CreateUser)
	user.Get("/:id", h.GetUserByID)
	user.Put("/:id", h.UpdateUserByID)
	user.Delete("/:id", h.DeleteUserByID)
}

// GetUserByID возвращает Юзера по ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	user, err := h.service.GetUserByID(c.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Error(c, errs.ErrNotFound, nil)
		}
		h.log.Error("Failed to get user", zap.Uint("id", uint(id)), zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, user, "")
}

// CreateUser создает и возвращает Юзера
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return errs.Error(c, errs.ErrTelegramUernameAlreadyExists, nil)
		}
		h.log.Error("Failed to create user", zap.Error(err))
		return errs.Error(c, errs.ErrInternal, nil)
	}

	return errs.Success(c, user, "")
}

// UpdateUserByID обновляет и возвращает Юзера
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if req.Name == nil && req.TelegramUsername == nil && req.Password == nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	user, err := h.service.UpdateUserByID(c.Context(), uint(id), req)
	if err != nil {
		switch {

		case errors.Is(err, gorm.ErrRecordNotFound):
			return errs.Error(c, errs.ErrNotFound, nil)

		case errors.Is(err, errors.New("no fields to update")):
			return errs.Error(c, errs.ErrEmptyBody, nil)

		case strings.Contains(err.Error(), "already exists"):
			return errs.Error(c, errs.ErrTelegramUernameAlreadyExists, nil)

		default:
			h.log.Error("Failed to update user", zap.Uint("id", uint(id)), zap.Error(err))
			return errs.Error(c, errs.ErrInternal, nil)
		}
	}
	return errs.Success(c, user, "")
}

// DeleteUserByID удаляет Юзера, если он не создал задачи или не назначен исполнителем
func (h *UserHandler) DeleteUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return errs.Error(c, errs.ErrInvalidID, nil)
	}

	err = h.service.DeleteUserByID(c.Context(), uint(id))
	if err != nil {
		switch {

		case errors.Is(err, gorm.ErrRecordNotFound):
			return errs.Error(c, errs.ErrNotFound, nil)

		case err.Error() == "User has associated tasks: Author":
			var createdTasks int64
			h.service.db.Model(&models.Task{}).Where("author_id = ?", id).Count(&createdTasks)

			return errs.Error(c, errs.ErrUserIsAuthorTasks, nil)

		case err.Error() == "User has associated tasks: Executor":
			var assignedTasks int64
			h.service.db.Model(&models.Task{}).Where("executor_id = ?", id).Count(&assignedTasks)

			return errs.Error(c, errs.ErrUserIsExecutorTasks, nil)

		default:
			h.log.Error("Failed to delete user", zap.Uint("id", uint(id)), zap.Error(err))
			return errs.Error(c, errs.ErrInternal, nil)
		}

	}
	return errs.Success(c, nil, "User deleted successfully")
}
