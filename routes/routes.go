// routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	RegisterKategoriRoutes(api)
	RegisterBarangRoutes(api)
	RegisterPeminjamanRoutes(api)
	RegisterLaporanRoutes(api)
}


