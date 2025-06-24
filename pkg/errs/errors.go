package errs

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Error(c *fiber.Ctx, error string, details interface{}) error {
	return c.Status(200).JSON(ErrorResponse{
		Success: false,
		Error:   error,
		Details: details,
	})
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(200).JSON(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
