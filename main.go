// @title Inventory Management API
// @version 1.0
// @description API untuk sistem manajemen inventory dengan authentication JWT
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token dengan format: Bearer {token}

package main

import (
	"inventory-backend/config"
	"inventory-backend/controllers"
	_ "inventory-backend/docs" // Import swagger docs
	"inventory-backend/middlewares"
	"inventory-backend/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using default values")
	}

	app := fiber.New()

	// Middleware
	middlewares.SetupMiddleware(app)

	// Connect DB
	config.ConnectDB()

	// ✅ Set DB ke controller setelah terkoneksi
	controllers.SetUserCollection(config.DB)
	controllers.SetKategoriCollection(config.DB)
	controllers.SetBarangCollection(config.DB)
	controllers.SetPeminjamanCollection(config.DB)
	controllers.SetLaporanCollection(config.DB)

	// Routes
	routes.SetupRoutes(app)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	startErr := app.Listen(":" + port)
	if startErr != nil {
		panic(startErr)
	}
}
