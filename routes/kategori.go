package routes

import (
	"inventory-backend/controllers"
	"inventory-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterKategoriRoutes(router fiber.Router) {
	kategori := router.Group("/kategori")
	
	// Public endpoints (semua user bisa akses)
	kategori.Get("/", middlewares.JWTMiddleware, controllers.GetAllKategori)
	kategori.Get("/:id", middlewares.JWTMiddleware, controllers.GetKategoriByID)
	
	// Protected endpoints (hanya admin yang bisa create/update/delete)
	kategori.Post("/", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.CreateKategori)
	kategori.Put("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.UpdateKategori)
	kategori.Delete("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.DeleteKategori)
}
