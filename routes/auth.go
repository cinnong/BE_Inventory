package routes

import (
	"inventory-backend/controllers"
	"inventory-backend/middlewares"
	"inventory-backend/validators"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Public routes (tidak perlu authentication)
	auth.Post("/register", validators.ValidateRegister, controllers.Register)
	auth.Post("/login", validators.ValidateLogin, controllers.Login)

	// Protected routes (perlu authentication)
	auth.Get("/profile", middlewares.JWTMiddleware, controllers.GetProfile)
}
