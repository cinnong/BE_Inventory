// routes/laporan.go
package routes

import (
	"inventory-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

// Route laporan untuk admin
func RegisterLaporanRoutes(router fiber.Router) {
	laporan := router.Group("/laporan")
	laporan.Get("/peminjaman", controllers.GetLaporanPeminjaman)
}
