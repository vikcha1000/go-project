package handlers

import (
    "github.com/gofiber/fiber/v2"
    "mine/internal/model"
    "mine/pkg/database"
)

func GetItems(c *fiber.Ctx) error {
    db := database.GetDB()
    items, err := models.GetAllItems(db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(items)
}

func CreateItem(c *fiber.Ctx) error {
    db := database.GetDB()
    
    var item models.Item
    if err := c.BodyParser(&item); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := models.CreateItem(db, &item); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(item)
}