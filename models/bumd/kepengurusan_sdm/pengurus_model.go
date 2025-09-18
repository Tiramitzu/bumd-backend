package kepengurusan_sdm

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PengurusModel struct {
	Id                  uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd              uuid.UUID `json:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	JabatanStruktur     int       `json:"jabatan_struktur" example:"1"`
	NamaPengurus        string    `json:"nama_pengurus" example:"Pengurus 1"`
	NIK                 string    `json:"nik" example:"1234567890"`
	Alamat              string    `json:"alamat" example:"Jl. Raya No. 1, Jakarta"`
	DeskripsiJabatan    string    `json:"deskripsi_jabatan" example:"Deskripsi Jabatan"`
	PendidikanAkhir     uuid.UUID `json:"pendidikan_akhir" example:"01994c79-6d4d-7e5e-9d30-3d28773ae539"`
	NamaPendidikanAkhir string    `json:"nama_pendidikan_akhir" example:"S3 (Doktor)"`
	TanggalMulaiJabatan time.Time `json:"tanggal_mulai_jabatan" example:"2021-01-01T00:00:00Z"`
	TanggalAkhirJabatan time.Time `json:"tanggal_akhir_jabatan" example:"2021-01-01T00:00:00Z"`
	File                string    `json:"file" example:"/path/to/file.pdf"`
}

type PengurusForm struct {
	JabatanStruktur     int                   `json:"jabatan_struktur" form:"jabatan_struktur" validate:"min=0,max=2" example:"1" default:"0"`
	NamaPengurus        string                `json:"nama_pengurus" form:"nama_pengurus" validate:"required" example:"Pengurus 1"`
	NIK                 string                `json:"nik" form:"nik" example:"1234567890"`
	Alamat              string                `json:"alamat" form:"alamat" example:"Jl. Raya No. 1, Jakarta"`
	DeskripsiJabatan    string                `json:"deskripsi_jabatan" form:"deskripsi_jabatan" example:"Deskripsi Jabatan"`
	PendidikanAkhir     uuid.UUID             `json:"pendidikan_akhir" form:"pendidikan_akhir" validate:"required" example:"01994c79-6d4d-7e5e-9d30-3d28773ae539"`
	TanggalMulaiJabatan string                `json:"tanggal_mulai_jabatan" form:"tanggal_mulai_jabatan" validate:"required,datetime=2006-01-02" example:"2021-01-01T00:00:00Z"`
	TanggalAkhirJabatan string                `json:"tanggal_akhir_jabatan" form:"tanggal_akhir_jabatan" validate:"required,datetime=2006-01-02" example:"2021-01-01T00:00:00Z"`
	File                *multipart.FileHeader `json:"file" form:"file"`
}
