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

// GetPeminjamanByID godoc
// @Summary Get peminjaman by ID
// @Description Mengambil data peminjaman berdasarkan ID
// @Tags Peminjaman
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Peminjaman ID"
// @Success 200 {object} models.Peminjaman "Data peminjaman"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Data peminjaman tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /peminjaman/{id} [get]
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

// GetAllPeminjaman godoc
// @Summary Get all peminjaman
// @Description Mengambil semua data peminjaman dengan opsi pencarian exact match nama peminjam
// @Tags Peminjaman
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Pencarian exact match berdasarkan nama peminjam (case-insensitive)"
// @Success 200 {array} models.Peminjaman "Daftar peminjaman"
// @Failure 500 {object} map[string]interface{} "Terjadi kesalahan server"
// @Router /peminjaman [get]
func GetAllPeminjaman(c *fiber.Ctx) error {
	search := c.Query("search")
	filter := bson.M{}

	// Jika parameter search ada, tambahkan filter nama_peminjam untuk exact match
	if search != "" {
		filter = bson.M{
			"nama_peminjam": bson.M{
				"$regex": primitive.Regex{
					Pattern: "^" + search + "$", // ^ untuk awal string, $ untuk akhir string (exact match)
					Options: "i",                // case-insensitive
				},
			},
		}
	}

	cursor, err := peminjamanCollection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data peminjaman"})
	}

	var peminjaman []models.Peminjaman
	if err := cursor.All(context.Background(), &peminjaman); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca data peminjaman"})
	}
	return c.JSON(peminjaman)
}

// CreatePeminjaman godoc
// @Summary Create new peminjaman
// @Description Membuat data peminjaman baru
// @Tags Peminjaman
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param peminjaman body models.Peminjaman true "Data peminjaman baru"
// @Success 201 {object} models.Peminjaman "Peminjaman berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Barang tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /peminjaman [post]
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

	// Hanya proses jika status dipinjam
	if data.Status == "dipinjam" {
		if data.Jumlah > barang.Stok {
			return c.Status(400).JSON(fiber.Map{"error": "Stok barang tidak mencukupi"})
		}
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

// UpdateStatusPeminjaman godoc
// @Summary Update status peminjaman
// @Description Mengupdate status peminjaman (dipinjam/dikembalikan)
// @Tags Peminjaman
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Peminjaman ID"
// @Param status body object{status=string} true "Status baru"
// @Success 200 {object} map[string]interface{} "Status berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Data tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /peminjaman/{id}/status [put]
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

	// Otomatisasi stok:
	if pinjam.Status != updateData.Status {
		if updateData.Status == "dipinjam" {
			// Validasi stok sebelum update
			var barang models.Barang
			err := barangCollectionPeminjaman.FindOne(context.Background(), bson.M{"_id": pinjam.BarangID}).Decode(&barang)
			if err != nil {
				return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
			}
			if pinjam.Jumlah > barang.Stok {
				return c.Status(400).JSON(fiber.Map{"error": "Stok barang tidak mencukupi"})
			}
			_, err = barangCollectionPeminjaman.UpdateOne(context.Background(),
				bson.M{"_id": pinjam.BarangID},
				bson.M{"$inc": bson.M{"stok": -pinjam.Jumlah}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		} else if pinjam.Status == "dipinjam" && updateData.Status == "dikembalikan" {
			_, err := barangCollectionPeminjaman.UpdateOne(context.Background(),
				bson.M{"_id": pinjam.BarangID},
				bson.M{"$inc": bson.M{"stok": pinjam.Jumlah}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
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

// DeletePeminjaman godoc
// @Summary Delete peminjaman
// @Description Menghapus data peminjaman berdasarkan ID
// @Tags Peminjaman
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Peminjaman ID"
// @Success 200 {object} map[string]interface{} "Data peminjaman berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Data peminjaman tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /peminjaman/{id} [delete]
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
