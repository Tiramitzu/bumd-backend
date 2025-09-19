package kepengurusan_sdm

import (
	"time"

	"github.com/google/uuid"
)

type PegawaiModel struct {
	Id             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd         uuid.UUID `json:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tahun          int       `json:"tahun" example:"2021"`
	StatusPegawai  int       `json:"status_pegawai" example:"1"`
	Pendidikan     uuid.UUID `json:"pendidikan" example:"01994c79-6d4d-7e5e-9d30-3d28773ae539"`
	NamaPendidikan string    `json:"nama_pendidikan" example:"S3 (Doktor)"`
	JumlahPegawai  int       `json:"jumlah_pegawai" example:"100"`
	CreatedAt      time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy      int64     `json:"created_by" example:"1"`
	UpdatedAt      time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy      int64     `json:"updated_by" example:"1"`
}

type PegawaiForm struct {
	Tahun         int       `json:"tahun" form:"tahun" validate:"required,min=2000,max=2099" example:"2021"`
	StatusPegawai int       `json:"status_pegawai" form:"status_pegawai" validate:"min=0,max=3" example:"1" default:"0"`
	Pendidikan    uuid.UUID `json:"pendidikan" form:"pendidikan" validate:"required" example:"01994c79-6d4d-7e5e-9d30-3d28773ae539"`
	JumlahPegawai int       `json:"jumlah_pegawai" form:"jumlah_pegawai" validate:"min=0" example:"100"`
}
