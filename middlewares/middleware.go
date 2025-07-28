package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://feinventory-production.up.railway.app, http://localhost:5173, https://beinventory-production.up.railway.app",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(logger.New())
}
