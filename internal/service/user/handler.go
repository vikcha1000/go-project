package user

import (
	"github.com/gofiber/fiber/v2"
)




func (h *UserHandler) SetupAPI(r fiber.Router) {
	user := r.Group("/user")
	//task.Post("/", taskHandler.CreateTask)
	user.Get("/:id", h.GetUserByID)
}