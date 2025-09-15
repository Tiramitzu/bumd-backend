package models

import "github.com/google/uuid"

type JenisDokumenModel struct {
	Id        uuid.UUID `json:"id" xml:"id"`
	Nama      string    `json:"nama" xml:"nama"`
	Deskripsi string    `json:"deskripsi" xml:"deskripsi"`
}

type JenisDokumenForm struct {
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}
