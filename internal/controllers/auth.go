package controllers

import (
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
	payload := &models.RegisterUserSchema{
		Username:     username,
		Email:        email,
		Password:     password,
		ProfileImage: pic,
	}

	authResponse, err := services.SignupUser(payload)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"authResponse": authResponse}})
}
