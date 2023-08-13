package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/golang-fiber/middleware"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/souvik150/golang-fiber/controllers"
	"github.com/souvik150/golang-fiber/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Route("/user", func(router fiber.Router) {
		router.Post("/login", controllers.LoginUser)
		router.Post("/signup", controllers.SignupUser)
		router.Post("/refresh", controllers.RefreshToken)
	})
	micro.Route("/notes", func(router fiber.Router) {
		router.Post("/", middleware.TokenValidation, controllers.CreateNoteHandler)
		router.Get("", middleware.TokenValidation, controllers.FindNotes)
	})

	micro.Route("/notes/:noteId", func(router fiber.Router) {
		router.Delete("", middleware.TokenValidation, controllers.DeleteNote)
		router.Get("", middleware.TokenValidation, controllers.FindNoteById)
		router.Patch("", middleware.TokenValidation, controllers.UpdateNote)
	})
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, and GORM",
		})
	})

	log.Fatal(app.Listen(":8000"))
}
