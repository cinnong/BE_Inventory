package routes

import (
	"inventory-backend/controllers"
	"inventory-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterPeminjamanRoutes(router fiber.Router) {
	peminjaman := router.Group("/peminjaman")
	
	// Semua user bisa lihat dan buat peminjaman
	peminjaman.Get("/", middlewares.JWTMiddleware, controllers.GetAllPeminjaman)
	peminjaman.Get("/:id", middlewares.JWTMiddleware, controllers.GetPeminjamanByID)
	peminjaman.Post("/", middlewares.JWTMiddleware, controllers.CreatePeminjaman)
	
	// Hanya admin yang bisa update status, jumlah, dan delete
	peminjaman.Put("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.UpdateStatusPeminjaman)
	peminjaman.Put("/:id/jumlah", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.UpdateJumlahPeminjaman)
	peminjaman.Delete("/:id", middlewares.JWTMiddleware, middlewares.RequireAdmin, controllers.DeletePeminjaman)
}
