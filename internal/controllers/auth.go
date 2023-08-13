package controllers

import (
	"github.com/gofiber/fiber/v2"
	database "github.com/souvik150/golang-fiber/internal/database"
	models "github.com/souvik150/golang-fiber/internal/models"
	token "github.com/souvik150/golang-fiber/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func SignupUser(c *fiber.Ctx) error {
	var payload models.RegisterUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to hash password"})
	}

	newUser := models.User{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	accessToken, err := token.GenerateAccessToken(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate access token"})
	}

	refreshToken, err := token.GenerateRefreshToken(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate refresh token"})
	}

	// Store the refresh token in the database
	refreshTokenEntry := models.RefreshToken{
		UserID: newUser.ID,
		Token:  refreshToken,
	}
	refreshTokenResult := database.DB.Create(&refreshTokenEntry)
	if refreshTokenResult.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": refreshTokenResult.Error.Error()})
	}

	authResponse := models.AuthResponse{
		UserID:       newUser.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": newUser, "authResponse": authResponse}})
}

func LoginUser(c *fiber.Ctx) error {
	var payload models.LoginUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := database.DB.Where("username = ?", payload.Username).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	accessToken, err := token.GenerateAccessToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate access token"})
	}

	refreshToken, err := token.GenerateRefreshToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate refresh token"})
	}

	// Store the refresh token in the database
	refreshTokenEntry := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}
	refreshTokenResult := database.DB.Create(&refreshTokenEntry)
	if refreshTokenResult.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": refreshTokenResult.Error.Error()})
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
}
