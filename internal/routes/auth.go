package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/golang-fiber/internal/controllers"
)

func AuthRoutes(group fiber.Router) {
	userGroup := group.Group("/user")
	userGroup.Post("/login", controllers.LoginUser)
	userGroup.Post("/signup", controllers.SignupUser)
	userGroup.Post("/refresh", controllers.RefreshToken)
}
