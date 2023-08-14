package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/souvik150/golang-fiber/internal/services"
)

func GetUserByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	user, err := services.GetUserByID(userID)
	if err != nil {
		// Handle error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Failed to get user details"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}
