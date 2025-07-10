// controllers/laporan.go
package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var dbRef *mongo.Database // Simpan referensi ke database

func SetLaporanCollection(db *mongo.Database) {
	dbRef = db
}

func GetLaporanPeminjaman(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Referensi koleksi
	peminjamanCollection := dbRef.Collection("peminjaman")

	// Pipeline aggregation untuk join ke barang dan kategori
	pipeline := mongo.Pipeline{
		// Join dengan koleksi barang
		{{Key: "$lookup", Value: bson.M{
			"from":         "barang",
			"localField":   "barang_id",
			"foreignField": "_id",
			"as":           "barang_info",
		}}},
		// Unwind array barang_info jadi objek tunggal
		{{Key: "$unwind", Value: "$barang_info"}},
		// Join dengan kategori
		{{Key: "$lookup", Value: bson.M{
			"from":         "kategori",
			"localField":   "barang_info.kategori_id",
			"foreignField": "_id",
			"as":           "kategori_info",
		}}},
		// Unwind kategori
		{{Key: "$unwind", Value: "$kategori_info"}},
	}

	cursor, err := peminjamanCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var hasil []bson.M
	if err := cursor.All(ctx, &hasil); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return data gabungan lengkap
	return c.JSON(hasil)
}
