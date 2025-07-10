package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Peminjaman struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NamaPeminjam   string             `json:"nama_peminjam" bson:"nama_peminjam"`
	EmailPeminjam  string             `json:"email_peminjam" bson:"email_peminjam"`
	TeleponPeminjam string            `json:"telepon_peminjam" bson:"telepon_peminjam"`
	BarangID       primitive.ObjectID `json:"barang_id" bson:"barang_id"`
	Jumlah         int                `json:"jumlah" bson:"jumlah"`
	TanggalPinjam  string             `json:"tanggal_pinjam" bson:"tanggal_pinjam"`
	Status         string             `json:"status" bson:"status"`
}

