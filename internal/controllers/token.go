package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/golang-fiber/internal/services"
)

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.FormValue("refreshToken")
	authResponse, err := services.RefreshAccessToken(refreshToken)
	if err != nil {
		// Handle error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to refresh access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
}
