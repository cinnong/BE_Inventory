package routes

import (
	"inventory-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterKategoriRoutes(router fiber.Router) {
	kategori := router.Group("/kategori")
	kategori.Get("/", controllers.GetAllKategori)
	kategori.Get("/:id", controllers.GetKategoriByID)
	kategori.Post("/", controllers.CreateKategori)
	kategori.Put("/:id", controllers.UpdateKategori)
	kategori.Delete("/:id", controllers.DeleteKategori)
}
