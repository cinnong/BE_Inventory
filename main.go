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
// @schemes http https

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

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
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

	// Start server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
