package dokumen

import "mime/multipart"

type PerdaPendirianModel struct {
	ID         int64   `json:"id" xml:"id"`
	Nomor      string  `json:"nomor" xml:"nomor"`
	Tanggal    string  `json:"tanggal" xml:"tanggal"`
	Keterangan string  `json:"keterangan" xml:"keterangan"`
	File       string  `json:"file" xml:"file"`
	ModalDasar float64 `json:"modal_dasar" xml:"modal_dasar"`
	IDBumd     int32   `json:"id_bumd" xml:"id_bumd"`
}

type PerdaPendirianForm struct {
	ID         *int64                `json:"id" xml:"id"`
	Nomor      string                `json:"nomor" xml:"nomor"`
	Tanggal    string                `json:"tanggal" xml:"tanggal"`
	Keterangan string                `json:"keterangan" xml:"keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file"`
	ModalDasar float64               `json:"modal_dasar" xml:"modal_dasar"`
}
