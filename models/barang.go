// models/barang.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Barang struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nama        string             `json:"nama" bson:"nama"`
	KategoriID  primitive.ObjectID `json:"kategori_id" bson:"kategori_id"`
	Stok        int                `json:"stok" bson:"stok"`
	TanggalBuat string             `json:"tanggal_buat" bson:"tanggal_buat"`
}
