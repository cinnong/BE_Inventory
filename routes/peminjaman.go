// routes/peminjaman.go
package routes

import (
	"inventory-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterPeminjamanRoutes(router fiber.Router) {
	peminjaman := router.Group("/peminjaman")
	peminjaman.Get("/", controllers.GetAllPeminjaman)
	peminjaman.Post("/", controllers.CreatePeminjaman)
	peminjaman.Put("/:id", controllers.UpdateStatusPeminjaman)
	peminjaman.Delete("/:id", controllers.DeletePeminjaman)
	peminjaman.Get("/:id", controllers.GetPeminjamanByID)
}
