package controllers

import (
	"context"
	"inventory-backend/models"
	"inventory-backend/validators"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var kategoriCollection *mongo.Collection

func SetKategoriCollection(db *mongo.Database) {
	kategoriCollection = db.Collection("kategori")
}

// GetAllKategori godoc
// @Summary Get all kategori
// @Description Mengambil semua data kategori
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Kategori "List semua kategori"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /kategori [get]
func GetAllKategori(c *fiber.Ctx) error {
	cursor, err := kategoriCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var kategori []models.Kategori
	if err := cursor.All(context.Background(), &kategori); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(kategori)
}

// GetKategoriByID godoc
// @Summary Get kategori by ID
// @Description Mengambil data kategori berdasarkan ID
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kategori ID"
// @Success 200 {object} models.Kategori "Data kategori"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Kategori tidak ditemukan"
// @Router /kategori/{id} [get]
func GetKategoriByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var kategori models.Kategori
	err = kategoriCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&kategori)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kategori tidak ditemukan"})
	}

	return c.JSON(kategori)
}

// CreateKategori godoc
// @Summary Create new kategori
// @Description Membuat data kategori baru
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param kategori body models.Kategori true "Data kategori baru"
// @Success 201 {object} models.Kategori "Kategori berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /kategori [post]
func CreateKategori(c *fiber.Ctx) error {
	var kategori models.Kategori
	if err := c.BodyParser(&kategori); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validators.ValidateKategori(kategori.Nama, kategori.Deskripsi); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	kategori.ID = primitive.NewObjectID()
	kategori.TanggalBuat = time.Now().Format("2006-01-02 15:04:05")

	_, err := kategoriCollection.InsertOne(context.Background(), kategori)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(kategori)
}

// UpdateKategori godoc
// @Summary Update kategori
// @Description Mengupdate data kategori berdasarkan ID
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kategori ID"
// @Param kategori body models.Kategori true "Data kategori yang akan diupdate"
// @Success 200 {object} map[string]interface{} "Kategori berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /kategori/{id} [put]
func UpdateKategori(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var data models.Kategori
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validators.ValidateKategori(data.Nama, data.Deskripsi); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"nama":         data.Nama,
			"deskripsi":    data.Deskripsi,
			"tanggal_buat": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	_, err = kategoriCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Kategori berhasil diupdate"})
}

// DeleteKategori godoc
// @Summary Delete kategori
// @Description Menghapus data kategori berdasarkan ID
// @Tags Kategori
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kategori ID"
// @Success 200 {object} map[string]interface{} "Kategori berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /kategori/{id} [delete]
func DeleteKategori(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	_, err = kategoriCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Kategori berhasil dihapus"})
}
