package controllers

import (
	"context"
	"inventory-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateJumlahPeminjaman mengubah jumlah peminjaman dan otomatis update stok barang
func UpdateJumlahPeminjaman(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var updateData struct {
		Jumlah int `json:"jumlah"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil data peminjaman lama
	var pinjam models.Peminjaman
	err = peminjamanCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&pinjam)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data peminjaman tidak ditemukan"})
	}

	if pinjam.Status != "dipinjam" {
		return c.Status(400).JSON(fiber.Map{"error": "Hanya peminjaman dengan status 'dipinjam' yang bisa diubah jumlahnya"})
	}

	if updateData.Jumlah <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Jumlah pinjam harus lebih dari 0"})
	}

	if updateData.Jumlah == pinjam.Jumlah {
		return c.Status(200).JSON(fiber.Map{"message": "Jumlah tidak berubah"})
	}

	// Ambil data barang
	var barang models.Barang
	err = barangCollectionPeminjaman.FindOne(context.Background(), bson.M{"_id": pinjam.BarangID}).Decode(&barang)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
	}

	diff := updateData.Jumlah - pinjam.Jumlah
	if diff > 0 {
		// Tambah jumlah pinjam, cek stok cukup
		if diff > barang.Stok {
			return c.Status(400).JSON(fiber.Map{"error": "Stok barang tidak mencukupi"})
		}
		_, err = barangCollectionPeminjaman.UpdateOne(context.Background(),
			bson.M{"_id": barang.ID},
			bson.M{"$inc": bson.M{"stok": -diff}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else if diff < 0 {
		// Kurangi jumlah pinjam, kembalikan stok
		_, err = barangCollectionPeminjaman.UpdateOne(context.Background(),
			bson.M{"_id": barang.ID},
			bson.M{"$inc": bson.M{"stok": -diff}}) // -diff karena diff negatif
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Update jumlah di peminjaman
	_, err = peminjamanCollection.UpdateOne(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"jumlah": updateData.Jumlah}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Jumlah peminjaman berhasil diubah"})
}
