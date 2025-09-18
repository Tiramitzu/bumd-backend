package dokumen

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type AktaNotarisModel struct {
	Id         uuid.UUID `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd     uuid.UUID `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nomor      string    `json:"nomor" xml:"nomor" example:"1234567890"`
	Notaris    string    `json:"notaris" xml:"notaris" example:"Notaris"`
	Tanggal    time.Time `json:"tanggal" xml:"tanggal" example:"2021-01-01T00:00:00Z"`
	Keterangan string    `json:"keterangan" xml:"keterangan" example:"Keterangan"`
	File       string    `json:"file" xml:"file" example:"/path/to/file.pdf"`
	CreatedAt  time.Time `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy  int64     `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt  time.Time `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy  int64     `json:"updated_by" xml:"updated_by" example:"1"`
}

type AktaNotarisForm struct {
	Nomor      string                `json:"nomor" xml:"nomor" form:"nomor" example:"1234567890"`
	Notaris    string                `json:"notaris" xml:"notaris" form:"notaris" example:"Notaris"`
	Tanggal    string                `json:"tanggal" xml:"tanggal" form:"tanggal" example:"2021-01-01T00:00:00Z"`
	Keterangan string                `json:"keterangan" xml:"keterangan" form:"keterangan" example:"Keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file" form:"file"`
}
