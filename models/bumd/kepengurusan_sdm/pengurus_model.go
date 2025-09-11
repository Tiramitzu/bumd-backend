package kepengurusan_sdm

import "mime/multipart"

type PengurusModel struct {
	ID                  int    `json:"id"`
	IDBumd              int    `json:"id_bumd"`
	JabatanStruktur     int    `json:"jabatan_struktur"`
	NamaPengurus        string `json:"nama_pengurus"`
	NIK                 string `json:"nik"`
	Alamat              string `json:"alamat"`
	DeskripsiJabatan    string `json:"deskripsi_jabatan"`
	PendidikanAkhir     int    `json:"pendidikan_akhir"`
	TanggalMulaiJabatan string `json:"tanggal_mulai_jabatan"`
	TanggalAkhirJabatan string `json:"tanggal_akhir_jabatan"`
	File                string `json:"file"`
}

type PengurusForm struct {
	JabatanStruktur     int                   `json:"jabatan_struktur" form:"jabatan_struktur"`
	NamaPengurus        string                `json:"nama_pengurus" form:"nama_pengurus"`
	NIK                 string                `json:"nik" form:"nik"`
	Alamat              string                `json:"alamat" form:"alamat"`
	DeskripsiJabatan    string                `json:"deskripsi_jabatan" form:"deskripsi_jabatan"`
	PendidikanAkhir     int                   `json:"pendidikan_akhir" form:"pendidikan_akhir"`
	TanggalMulaiJabatan string                `json:"tanggal_mulai_jabatan" form:"tanggal_mulai_jabatan"`
	TanggalAkhirJabatan string                `json:"tanggal_akhir_jabatan" form:"tanggal_akhir_jabatan"`
	File                *multipart.FileHeader `json:"file" form:"file"`
}
