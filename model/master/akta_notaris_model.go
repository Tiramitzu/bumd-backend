package models

type AktaNotarisModel struct {
	ID           int     `json:"id" xml:"id"`
	IDBumd       int     `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
}

type AktaNotarisForm struct {
	IDBumd       int     `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
}
