package models

import "github.com/google/uuid"

type JenisLaporanModel struct {
	Id          int64     `json:"id" xml:"id"`
	BentukUsaha uuid.UUID `json:"bentuk_usaha" xml:"bentuk_usaha"`
	Kode        string    `json:"kode" xml:"kode"`
	Uraian      string    `json:"uraian" xml:"uraian"`
	Keterangan  string    `json:"keterangan" xml:"keterangan"`
	Level       int       `json:"level" xml:"level"`
	ParentId    int       `json:"parent_id" xml:"parent_id"`
}
