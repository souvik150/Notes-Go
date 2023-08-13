package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	database "github.com/souvik150/golang-fiber/internal/database"
	"github.com/souvik150/golang-fiber/internal/middleware"
	models "github.com/souvik150/golang-fiber/internal/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateNoteHandler(c *fiber.Ctx) error {
	var payload *models.CreateNoteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Get the user ID from the authenticated user (adjust according to your authentication logic)
	userID := c.Locals("userID").(uuid.UUID)

	errors := middleware.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newNote := models.Note{
		UserID:    userID, // Set the user ID
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Published: payload.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := database.DB.Create(&newNote)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exists, please use another title"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": newNote}})
}

func FindNotes(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var notes []models.Note
	results := database.DB.Limit(intLimit).Offset(offset).Find(&notes)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(notes), "notes": notes})
}

func UpdateNote(c *fiber.Ctx) error {
	noteID := c.Params("noteId")

	var payload *models.UpdateNoteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	userID := c.Locals("userID").(uuid.UUID)

	var note models.Note
	result := database.DB.First(&note, "id = ? AND user_id = ?", noteID, userID)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that ID exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Category != "" {
		updates["category"] = payload.Category
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}

	if payload.Published != nil {
		updates["published"] = payload.Published
	}

	updates["updated_at"] = time.Now()

	database.DB.Model(&note).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": note}})
}

func FindNoteById(c *fiber.Ctx) error {
	noteId := c.Params("noteId")

	var note models.Note
	result := database.DB.First(&note, "id = ?", noteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": note}})
}

func DeleteNote(c *fiber.Ctx) error {
	noteId := c.Params("noteId")

	result := database.DB.Delete(&models.Note{}, "id = ?", noteId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
