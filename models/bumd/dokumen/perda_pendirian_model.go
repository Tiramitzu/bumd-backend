package dokumen

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PerdaPendirianModel struct {
	Id         uuid.UUID `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd     uuid.UUID `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nomor      string    `json:"nomor" xml:"nomor" example:"1234567890"`
	Tanggal    time.Time `json:"tanggal" xml:"tanggal" example:"2021-01-01T00:00:00Z"`
	Keterangan string    `json:"keterangan" xml:"keterangan" example:"Keterangan"`
	File       string    `json:"file" xml:"file" example:"/path/to/file.pdf"`
	ModalDasar float64   `json:"modal_dasar" xml:"modal_dasar" example:"1000000"`
	CreatedAt  time.Time `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy  int64     `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt  time.Time `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy  int64     `json:"updated_by" xml:"updated_by" example:"1"`
}

type PerdaPendirianForm struct {
	Nomor      string                `json:"nomor_perda" xml:"nomor_perda" form:"nomor_perda" validate:"required" example:"1234567890"`
	Tanggal    string                `json:"tanggal_perda" xml:"tanggal_perda" form:"tanggal_perda" validate:"required,datetime=2006-01-02" example:"2021-01-01"`
	Keterangan string                `json:"keterangan" xml:"keterangan" form:"keterangan" example:"Keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file" form:"file"`
	ModalDasar string                `json:"modal_dasar" xml:"modal_dasar" form:"modal_dasar" validate:"min=0" example:"1000000"`
}
