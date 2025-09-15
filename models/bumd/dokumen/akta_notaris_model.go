package dokumen

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type AktaNotarisModel struct {
	Id         uuid.UUID `json:"id" xml:"id"`
	IdBumd     uuid.UUID `json:"id_bumd" xml:"id_bumd"`
	Nomor      string    `json:"nomor" xml:"nomor"`
	Notaris    string    `json:"notaris" xml:"notaris"`
	Tanggal    time.Time `json:"tanggal" xml:"tanggal"`
	Keterangan string    `json:"keterangan" xml:"keterangan"`
	File       string    `json:"file" xml:"file"`
}

type AktaNotarisForm struct {
	Nomor      string                `json:"nomor" xml:"nomor" form:"nomor"`
	Notaris    string                `json:"notaris" xml:"notaris" form:"notaris"`
	Tanggal    string                `json:"tanggal" xml:"tanggal" form:"tanggal"`
	Keterangan string                `json:"keterangan" xml:"keterangan" form:"keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file" form:"file"`
}
