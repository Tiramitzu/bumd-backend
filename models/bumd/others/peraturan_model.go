package others

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PeraturanModel struct {
	Id                  uuid.UUID  `json:"id" xml:"id"`
	IdBumd              uuid.UUID  `json:"id_bumd" xml:"id_bumd"`
	Nomor               string     `json:"nomor" xml:"nomor"`
	TanggalBerlaku      *time.Time `json:"tanggal_berlaku" xml:"tanggal_berlaku"`
	KeteranganPeraturan string     `json:"keterangan_peraturan" xml:"keterangan_peraturan"`
	FilePeraturan       string     `json:"file_peraturan" xml:"file_peraturan"`
	JenisPeraturan      int32      `json:"jenis_peraturan" xml:"jenis_peraturan"`
	NamaJenisPeraturan  string     `json:"nama_jenis_peraturan" xml:"nama_jenis_peraturan"`
}

type PeraturanForm struct {
	Nomor               string                `json:"nomor" xml:"nomor" form:"nomor"`
	TanggalBerlaku      *string               `json:"tanggal_berlaku" xml:"tanggal_berlaku" form:"tanggal_berlaku"`
	KeteranganPeraturan string                `json:"keterangan_peraturan" xml:"keterangan_peraturan" form:"keterangan_peraturan"`
	FilePeraturan       *multipart.FileHeader `json:"file_peraturan" xml:"file_peraturan" form:"file_peraturan"`
	JenisPeraturan      int32                 `json:"jenis_peraturan" xml:"jenis_peraturan" form:"jenis_peraturan"`
}
