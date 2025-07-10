package main

import (
	"inventory-backend/config"
	"inventory-backend/controllers"
	"inventory-backend/middlewares"
	"inventory-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Middleware
	middlewares.SetupMiddleware(app)

	// Connect DB
	config.ConnectDB()

	// ✅ Set DB ke controller setelah terkoneksi
	controllers.SetKategoriCollection(config.DB)
	controllers.SetBarangCollection(config.DB)
	controllers.SetPeminjamanCollection(config.DB)
	controllers.SetLaporanCollection(config.DB)

	// Routes
	routes.SetupRoutes(app)

	// Start server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
