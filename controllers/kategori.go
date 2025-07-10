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
