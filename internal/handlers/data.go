package handlers

import (
	"mine/internal/model"

	"github.com/gofiber/fiber/v2"
)

func GetData(c *fiber.Ctx) error {
	data, err := model.GetAllData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

func CreateData(c *fiber.Ctx) error {
	var input model.DataInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	result, err := model.CreateData(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
