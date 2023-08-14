package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/golang-fiber/internal/models"
	"github.com/souvik150/golang-fiber/internal/services"
	"github.com/souvik150/golang-fiber/internal/utils"
	"log"
	"mime/multipart"
)

func SignupUser(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Failed to process form data"})
	}

	// Get user data from the form
	username := form.Value["username"][0]
	email := form.Value["email"][0]
	password := form.Value["password"][0]

	// Handle profile picture upload
	files := form.File["profilePic"]
	pic := ""

	for _, file := range files {
		fileHeader := file

		f, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer func(f multipart.File) {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(f)

		uploadedURL, err := utils.UploadFile(f, fileHeader)
		pic = uploadedURL
	}
	// Create a payload for user registration
	fmt.Println(pic)
	payload := &models.RegisterUserSchema{
		Username:     username,
		Email:        email,
		Password:     password,
		ProfileImage: pic,
	}

	err = services.SignupUser(payload)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Account created successfully. Please verify your account."})
}

func VerifyOTP(c *fiber.Ctx) error {
	var request models.VerifyOTPRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid request body"})
	}

	// Fetch user by email
	user, err := services.GetUserByEmail(request.Email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	err = services.VerifyOTP(user.ID, request.OTP)
	if err != nil {
		log.Println("Error verifying OTP:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid OTP"})
	}

	// Generate auth tokens for the verified user
	authResponse, err := services.GenerateAuthTokens(&user)
	if err != nil {
		log.Println("Error generating auth tokens:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Failed to generate auth tokens"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
}

func ResendOTP(c *fiber.Ctx) error {
	var request models.ResendOTPRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid request body"})
	}

	user, err := services.GetUserByEmail(request.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Failed to fetch user data"})
	}

	if user.Verified {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User is already verified"})
	}

	err = services.ResendOTP(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Failed to resend OTP"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "OTP resent successfully"})
}

func LoginUser(c *fiber.Ctx) error {
	var payload models.LoginUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	authResponse, err := services.LoginUser(&payload)
	if err != nil {
		// Handle error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid credentials"})
	}

	if !authResponse.Verified {
		// User is not verified
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Account not verified"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
}
