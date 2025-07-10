package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Kategori struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nama        string             `json:"nama" bson:"nama"`
	Deskripsi   string             `json:"deskripsi" bson:"deskripsi"`
	TanggalBuat string             `json:"tanggal_buat" bson:"tanggal_buat"`
}
