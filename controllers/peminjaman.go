// controllers/peminjaman.go
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

var peminjamanCollection *mongo.Collection
var barangCollectionPeminjaman *mongo.Collection

func SetPeminjamanCollection(db *mongo.Database) {
	peminjamanCollection = db.Collection("peminjaman")
	barangCollectionPeminjaman = db.Collection("barang")
}


func GetPeminjamanByID(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var peminjaman models.Peminjaman
	err = peminjamanCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&peminjaman)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Data peminjaman tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(peminjaman)
}

func GetAllPeminjaman(c *fiber.Ctx) error {
	cursor, err := peminjamanCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var peminjaman []models.Peminjaman
	if err := cursor.All(context.Background(), &peminjaman); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(peminjaman)
}

func CreatePeminjaman(c *fiber.Ctx) error {
	var data models.Peminjaman
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validators.ValidatePeminjaman(data.NamaPeminjam, data.EmailPeminjam, data.TeleponPeminjam, data.Jumlah, data.Status); err != nil {
	return c.Status(400).JSON(fiber.Map{"error": err.Error()})
}

	// Cek barang
	var barang models.Barang
	err := barangCollectionPeminjaman.FindOne(context.Background(), bson.M{"_id": data.BarangID}).Decode(&barang)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
	}

	// Cek stok
	if data.Jumlah > barang.Stok {
		return c.Status(400).JSON(fiber.Map{"error": "Stok barang tidak mencukupi"})
	}

	// Kurangi stok jika status dipinjam
	if data.Status == "dipinjam" {
		_, err = barangCollectionPeminjaman.UpdateOne(context.Background(),
			bson.M{"_id": barang.ID},
			bson.M{"$inc": bson.M{"stok": -data.Jumlah}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	data.ID = primitive.NewObjectID()
	data.TanggalPinjam = time.Now().Format("2006-01-02 15:04:05")

	_, err = peminjamanCollection.InsertOne(context.Background(), data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(data)
}

func UpdateStatusPeminjaman(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var updateData struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil data peminjaman dulu
	var pinjam models.Peminjaman
	err = peminjamanCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&pinjam)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	// Jika status dari dipinjam ke dikembalikan, tambahkan stok barang
	if pinjam.Status == "dipinjam" && updateData.Status == "dikembalikan" {
		_, err := barangCollection.UpdateOne(context.Background(),
			bson.M{"_id": pinjam.BarangID},
			bson.M{"$inc": bson.M{"stok": pinjam.Jumlah}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Update status
	_, err = peminjamanCollection.UpdateOne(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": updateData.Status}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Status berhasil diperbarui"})
}

func DeletePeminjaman(c *fiber.Ctx) error {
    id, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    // Cari data peminjaman yang akan dihapus
    var peminjaman models.Peminjaman
    err = peminjamanCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&peminjaman)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data peminjaman tidak ditemukan"})
    }

    // Jika status masih "dipinjam", kembalikan stok barang
    if peminjaman.Status == "dipinjam" {
        _, err = barangCollectionPeminjaman.UpdateOne(
            context.Background(),
            bson.M{"_id": peminjaman.BarangID},
            bson.M{"$inc": bson.M{"stok": peminjaman.Jumlah}},
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Gagal mengembalikan stok barang"})
        }
    }

    // Hapus data peminjaman
    _, err = peminjamanCollection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data peminjaman"})
    }

    return c.Status(200).JSON(fiber.Map{"message": "Data peminjaman berhasil dihapus"})
}
