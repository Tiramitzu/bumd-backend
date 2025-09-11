package dokumen

import (
	"mime/multipart"
	"time"
)

type AktaNotarisModel struct {
	ID         int64     `json:"id" xml:"id"`
	Nomor      string    `json:"nomor" xml:"nomor"`
	Notaris    string    `json:"notaris" xml:"notaris"`
	Tanggal    time.Time `json:"tanggal" xml:"tanggal"`
	Keterangan string    `json:"keterangan" xml:"keterangan"`
	File       string    `json:"file" xml:"file"`
	IDBumd     int32     `json:"id_bumd" xml:"id_bumd"`
}

type AktaNotarisForm struct {
	Nomor      string                `json:"nomor" xml:"nomor" form:"nomor"`
	Notaris    string                `json:"notaris" xml:"notaris" form:"notaris"`
	Tanggal    string                `json:"tanggal" xml:"tanggal" form:"tanggal"`
	Keterangan string                `json:"keterangan" xml:"keterangan" form:"keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file" form:"file"`
}
