// routes/barang.go
package routes

import (
	"inventory-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterBarangRoutes(router fiber.Router) {
	barang := router.Group("/barang")
	barang.Get("/", controllers.GetAllBarang)
	barang.Get("/:id", controllers.GetBarangByID)
	barang.Post("/", controllers.CreateBarang)
	barang.Put("/:id", controllers.UpdateBarang)
	barang.Delete("/:id", controllers.DeleteBarang)
}
