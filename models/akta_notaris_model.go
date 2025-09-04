package models

import "mime/multipart"

type AktaNotarisModel struct {
	ID           int64   `json:"id" xml:"id"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	IDBumd       int32   `json:"id_bumd" xml:"id_bumd"`
}

type AktaNotarisForm struct {
	NomorPerda   string                `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string                `json:"tanggal_perda" xml:"tanggal_perda"`
	Keterangan   string                `json:"keterangan" xml:"keterangan"`
	File         *multipart.FileHeader `json:"file" xml:"file"`
	ModalDasar   float64               `json:"modal_dasar" xml:"modal_dasar"`
	IDBumd       int32                 `json:"id_bumd" xml:"id_bumd"`
}
