package others

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PeraturanModel struct {
	Id                  uuid.UUID  `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd              uuid.UUID  `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nomor               string     `json:"nomor" xml:"nomor" example:"1234567890"`
	TanggalBerlaku      *time.Time `json:"tanggal_berlaku" xml:"tanggal_berlaku" example:"2021-01-01T00:00:00Z"`
	KeteranganPeraturan string     `json:"keterangan_peraturan" xml:"keterangan_peraturan" example:"Keterangan Peraturan"`
	FilePeraturan       string     `json:"file_peraturan" xml:"file_peraturan" example:"/path/to/file.pdf"`
	JenisPeraturan      uuid.UUID  `json:"jenis_peraturan" xml:"jenis_peraturan" example:"01994c04-699d-75e0-a288-13980f8c854d"`
	NamaJenisPeraturan  string     `json:"nama_jenis_peraturan" xml:"nama_jenis_peraturan" example:"SOP"`
	CreatedAt           time.Time  `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy           int64      `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt           time.Time  `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy           int64      `json:"updated_by" xml:"updated_by" example:"1"`
}

type PeraturanForm struct {
	Nomor               string                `json:"nomor" xml:"nomor" form:"nomor" validate:"required" example:"1234567890"`
	TanggalBerlaku      *string               `json:"tanggal_berlaku" xml:"tanggal_berlaku" form:"tanggal_berlaku" validate:"required,datetime=2006-01-02" example:"2021-01-01T00:00:00Z"`
	KeteranganPeraturan string                `json:"keterangan_peraturan" xml:"keterangan_peraturan" form:"keterangan_peraturan" example:"Keterangan Peraturan"`
	FilePeraturan       *multipart.FileHeader `json:"file_peraturan" xml:"file_peraturan" form:"file_peraturan"`
	JenisPeraturan      uuid.UUID             `json:"jenis_peraturan" xml:"jenis_peraturan" form:"jenis_peraturan" validate:"required" example:"01994c04-699d-75e0-a288-13980f8c854d"`
}
