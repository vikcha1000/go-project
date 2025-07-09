package login

import (
	"mine/pkg/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoginHandler struct {
	service  *LoginService
	validate *validator.Validate
	log      *zap.Logger
}

func NewLoginHandler(service *LoginService, log *zap.Logger) *LoginHandler {
	return &LoginHandler{
		service:  service,
		validate: validator.New(),
		log:      log,
	}
}

func (h *LoginHandler) SetupAPI(r fiber.Router) {
	login := r.Group("/login")
	login.Post("/", h.Login)
}

// Login проверяет совпадение введенного пароля
func (h *LoginHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.Error(c, errs.ErrInvalidBody, nil)
	}

	login, err := h.service.Login(c.Context(), req)
	if err != nil {
		return errs.Error(c, errs.ErrInvalidCreditails, nil)
	}

	return errs.Success(c, login, "")
}
