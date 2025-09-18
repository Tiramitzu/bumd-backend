package keuangan

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type LaporanKeuanganModel struct {
	Id                   uuid.UUID `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd               uuid.UUID `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdJenisLaporan       int       `json:"id_jenis_laporan" xml:"id_jenis_laporan" example:"2"`
	NamaJenisLaporan     string    `json:"nama_jenis_laporan" xml:"nama_jenis_laporan" example:"Cakupan Layanan"`
	KodeJenisLaporan     string    `json:"kode_jenis_laporan" xml:"kode_jenis_laporan" example:"1.1"`
	IdJenisLaporanItem   int       `json:"id_jenis_laporan_item" xml:"id_jenis_laporan_item" example:"3"`
	NamaJenisLaporanItem string    `json:"nama_jenis_laporan_item" xml:"nama_jenis_laporan_item" example:"Penduduk Wilayah Administrasi (Jiwa)"`
	KodeJenisLaporanItem string    `json:"kode_jenis_laporan_item" xml:"kode_jenis_laporan_item" example:"1.1.01"`
	Tahun                int       `json:"tahun" xml:"tahun" example:"2021"`
	Jumlah               float64   `json:"jumlah" xml:"jumlah" example:"1000000"`
	File                 string    `json:"file" xml:"file" example:"/path/to/file.pdf"`
	CreatedAt            time.Time `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy            int       `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt            time.Time `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy            int       `json:"updated_by" xml:"updated_by" example:"1"`
}

type LaporanKeuanganForm struct {
	IdJenisLaporan     int                   `json:"id_jenis_laporan" xml:"id_jenis_laporan" form:"id_jenis_laporan" validate:"required" example:"2"`
	IdJenisLaporanItem int                   `json:"id_jenis_laporan_item" xml:"id_jenis_laporan_item" form:"id_jenis_laporan_item" validate:"required" example:"3"`
	Tahun              int                   `json:"tahun" xml:"tahun" form:"tahun" validate:"required,min=2000,max=2099" example:"2021"`
	Jumlah             float64               `json:"jumlah" xml:"jumlah" form:"jumlah" validate:"required,min=0" example:"1000000"`
	File               *multipart.FileHeader `json:"file" xml:"file" form:"file"`
}
