// controllers/barang.go
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

var barangCollection *mongo.Collection


func SetBarangCollection(db *mongo.Database) {
	barangCollection = db.Collection("barang")
}


func GetAllBarang(c *fiber.Ctx) error {
	cursor, err := barangCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var barang []models.Barang
	if err := cursor.All(context.Background(), &barang); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(barang)
}

func GetBarangByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var barang models.Barang
	err = barangCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&barang)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
	}

	return c.JSON(barang)
}

func CreateBarang(c *fiber.Ctx) error {
	var barang models.Barang
	if err := c.BodyParser(&barang); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Validasi input
	// Pastikan KategoriID valid
	if err := validators.ValidateBarang(barang.Nama, barang.Stok); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Cek apakah kategori id ada
	var kategori models.Kategori
	err := kategoriCollection.FindOne(context.Background(), bson.M{"_id": barang.KategoriID}).Decode(&kategori)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Kategori tidak ditemukan"})
	}
	
	// INSERT DATA BARU
	barang.ID = primitive.NewObjectID()
	barang.TanggalBuat = time.Now().Format("2006-01-02 15:04:05")

	_, err = barangCollection.InsertOne(context.Background(), barang)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(barang)
}

func UpdateBarang(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var data models.Barang
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validators.ValidateBarang(data.Nama, data.Stok); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Pastikan KategoriID valid
	var kategori models.Kategori
	err = kategoriCollection.FindOne(context.Background(), bson.M{"_id": data.KategoriID}).Decode(&kategori)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Kategori tidak ditemukan"})
	}

	update := bson.M{
		"$set": bson.M{
			"nama":         data.Nama,
			"kategori_id":  data.KategoriID,
			"stok":         data.Stok,
			"tanggal_buat": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	_, err = barangCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Barang berhasil diupdate"})
}

func DeleteBarang(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	_, err = barangCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Barang berhasil dihapus"})
}
