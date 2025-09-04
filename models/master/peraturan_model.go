package models

type PeraturanModel struct {
	ID             int64  `json:"id" xml:"id"`
	IDBumd         int32  `json:"id_bumd" xml:"id_bumd"`
	Nomor          string `json:"nomor" xml:"nomor"`
	Jenis          string `json:"jenis" xml:"jenis"`
	TanggalBerlaku string `json:"tanggal_berlaku" xml:"tanggal_berlaku"`
	Keterangan     string `json:"keterangan" xml:"keterangan"`
	File           string `json:"file" xml:"file"`
}

type PeraturanForm struct {
	IDBumd         int32  `json:"id_bumd" xml:"id_bumd"`
	Nomor          string `json:"nomor" xml:"nomor"`
	Jenis          string `json:"jenis" xml:"jenis"`
	TanggalBerlaku string `json:"tanggal_berlaku" xml:"tanggal_berlaku"`
	Keterangan     string `json:"keterangan" xml:"keterangan"`
	File           string `json:"file" xml:"file"`
}
