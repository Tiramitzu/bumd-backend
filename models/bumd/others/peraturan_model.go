package others

import (
	"mime/multipart"
	"time"
)

type PeraturanModel struct {
	ID                  string     `json:"id" xml:"id"`
	Nomor               string     `json:"nomor" xml:"nomor"`
	TanggalBerlaku      *time.Time `json:"tanggal_berlaku" xml:"tanggal_berlaku"`
	KeteranganPeraturan string     `json:"keterangan_peraturan" xml:"keterangan_peraturan"`
	FilePeraturan       string     `json:"file_peraturan" xml:"file_peraturan"`
	IDBumd              int32      `json:"id_bumd" xml:"id_bumd"`
	JenisPeraturan      int32      `json:"jenis_peraturan" xml:"jenis_peraturan"`
	NamaJenisPeraturan  string     `json:"nama_jenis_peraturan" xml:"nama_jenis_peraturan"`
}

type PeraturanForm struct {
	Nomor               string                `json:"nomor" xml:"nomor" form:"nomor"`
	TanggalBerlaku      *string               `json:"tanggal_berlaku" xml:"tanggal_berlaku" form:"tanggal_berlaku"`
	KeteranganPeraturan string                `json:"keterangan_peraturan" xml:"keterangan_peraturan" form:"keterangan_peraturan"`
	FilePeraturan       *multipart.FileHeader `json:"file_peraturan" xml:"file_peraturan" form:"file_peraturan"`
	IDBumd              int32                 `json:"id_bumd" xml:"id_bumd" form:"id_bumd"`
	JenisPeraturan      int32                 `json:"jenis_peraturan" xml:"jenis_peraturan" form:"jenis_peraturan"`
}
