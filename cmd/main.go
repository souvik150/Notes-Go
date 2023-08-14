package main

import (
	"github.com/gofiber/fiber/v2"
	database "github.com/souvik150/golang-fiber/internal/database"
	"github.com/souvik150/golang-fiber/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/souvik150/golang-fiber/config"
)

func main() {
	app := fiber.New()

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	database.ConnectDB(&config)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	apiGroup := app.Group("/v1")

	routes.AuthRoutes(apiGroup)
	//routes.UserRoutes(apiGroup)
	routes.NotesRoutes(apiGroup)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})

	app.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Notes App",
		})
	})

	log.Fatal(app.Listen(config.Port))
}
