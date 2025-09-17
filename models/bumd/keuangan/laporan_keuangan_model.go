package keuangan

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type LaporanKeuanganModel struct {
	Id                 uuid.UUID `json:"id" xml:"id"`
	IdBumd             uuid.UUID `json:"id_bumd" xml:"id_bumd"`
	IdJenisLaporan     uuid.UUID `json:"id_jenis_laporan" xml:"id_jenis_laporan"`
	IdJenisLaporanItem uuid.UUID `json:"id_jenis_laporan_item" xml:"id_jenis_laporan_item"`
	Tahun              int       `json:"tahun" xml:"tahun"`
	Jumlah             float64   `json:"jumlah" xml:"jumlah"`
	File               string    `json:"file" xml:"file"`
	CreatedAt          time.Time `json:"created_at" xml:"created_at"`
	CreatedBy          int       `json:"created_by" xml:"created_by"`
	UpdatedAt          time.Time `json:"updated_at" xml:"updated_at"`
	UpdatedBy          int       `json:"updated_by" xml:"updated_by"`
}

type LaporanKeuanganForm struct {
	IdBumd             uuid.UUID             `json:"id_bumd" xml:"id_bumd"`
	IdJenisLaporan     int                   `json:"id_jenis_laporan" xml:"id_jenis_laporan"`
	IdJenisLaporanItem int                   `json:"id_jenis_laporan_item" xml:"id_jenis_laporan_item"`
	Tahun              int                   `json:"tahun" xml:"tahun"`
	Jumlah             float64               `json:"jumlah" xml:"jumlah"`
	File               *multipart.FileHeader `json:"file" xml:"file"`
}
