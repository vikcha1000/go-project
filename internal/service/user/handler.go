package user

import (
	"errors"

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
	if err != nil {
		h.log.Warn("Invalid user ID")
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный id")
	}
	user, err := h.service.GetUserByID(c.Context(), uint(id))
	if err != nil {
		h.log.Warn("User not found", zap.Uint("id", uint(id)))
		return fiber.NewError(fiber.StatusNotFound, "Юзер не найден")
	}
	return c.JSON(*user)
}

// CreateUser создает и возвращает Юзера
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Некорректное тело запроса")
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to create user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return c.JSON(*user)
}

// UpdateUserByID обновляет и возвращает Юзера
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.log.Warn("Invalid user ID")
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный id")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Failed to parse request", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Некорректное тело запроса")
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warn("Validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.service.UpdateUserByID(c.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Warn("User not found", zap.Uint("id", uint(id)))
			return fiber.NewError(fiber.StatusNotFound, "Юзер не найден")
		}
		if errors.Is(err, errors.New("no fields to update")) {
			return fiber.NewError(fiber.StatusBadRequest, "Нет полей для обновления")
		}
		h.log.Error("Failed to update user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return c.JSON(*user)
}

// DeleteUserByID удаляет Юзера, если он не создал задачи или не назначен исполнителем
func (h *UserHandler) DeleteUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.log.Warn("Invalid user ID")
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный id")
	}
	err = h.service.DeleteUserByID(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			h.log.Warn("User not found", zap.Uint("id", uint(id)))
			return fiber.NewError(fiber.StatusNotFound, "Юзер не найден")
		case err.Error() == "User has associated tasks: Author":
			h.log.Warn("User has associated tasks", zap.Uint("id", uint(id)))
			return fiber.NewError(fiber.StatusConflict, "Не возможно удалить Юзера, который создал задачи")
		case err.Error() == "User has associated tasks: Executor":
			h.log.Warn("User has associated tasks", zap.Uint("id", uint(id)))
			return fiber.NewError(fiber.StatusConflict, "Не возможно удалить Юзера, у которого есть назначенные задачи")
		default:
			h.log.Error("Failed to delete user", zap.Error(err))
			return fiber.ErrInternalServerError
		}
	}
	return c.SendStatus(fiber.StatusNoContent)

}
