package models

import "github.com/google/uuid"

type BentukUsahaModel struct {
	Id        uuid.UUID `json:"id" xml:"id"`
	Nama      string    `json:"nama" xml:"nama"`
	Deskripsi string    `json:"deskripsi" xml:"deskripsi"`
}

type BentukUsahaForm struct {
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}
