package models

type PerdaModel struct {
	ID           int64   `json:"id" xml:"id"`
	IDBumd       int32   `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
}

type PerdaForm struct {
	IDBumd       int32   `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
}
