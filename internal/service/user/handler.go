package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
}

// GetUserByID возвращает Юзера по ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}
	user, err := h.service.GetUserByID(c.Context(), uint(id))
	if err != nil {
		h.log.Warn("User not found", zap.Uint("id", uint(id)))
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	return c.JSON(*user)
}

// CreateUser создает и возвращает Юзера 
func (s *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		s.log.Warn("Failed to parse request", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

		if err := s.validate.Struct(req); err != nil {
		s.log.Warn("Validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

		user, err := s.service.CreateUser(c.Context(), req)
	if err != nil {
		s.log.Error("Failed to create user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return c.JSON(*user)
}