package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://feinventory-production.up.railway.app, http://localhost:5173, http://localhost:5174", // contoh domain react kamu di railway
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(logger.New())
}
