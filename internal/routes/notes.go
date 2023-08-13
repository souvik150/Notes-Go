package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/golang-fiber/internal/controllers"
	"github.com/souvik150/golang-fiber/internal/middleware"
)

func NotesRoutes(group fiber.Router) {
	notesGroup := group.Group("/notes")
	notesGroup.Use(middleware.TokenValidation)
	notesGroup.Post("/", controllers.CreateNoteHandler)
	notesGroup.Get("", controllers.FindNotes)

	notesGroup.Route("/:noteId", func(router fiber.Router) {
		router.Delete("", controllers.DeleteNote)
		router.Get("", controllers.FindNoteById)
		router.Patch("", controllers.UpdateNote)
	})
}
