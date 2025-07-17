package routes

import (
	"inventory-backend/controllers"
	"inventory-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterBarangRoutes(router fiber.Router) {
	barang := router.Group("/barang")
	
	// Public endpoints (semua user bisa akses)
	barang.Get("/", middlewares.JWTMiddleware, controllers.GetAllBarang)
	barang.Get("/:id", middlewares.JWTMiddleware, controllers.GetBarangByID)
	
	// Protected endpoints (hanya admin yang bisa create/update/delete)
	barang.Post("/", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.CreateBarang)
	barang.Put("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.UpdateBarang)
	barang.Delete("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.DeleteBarang)
}
