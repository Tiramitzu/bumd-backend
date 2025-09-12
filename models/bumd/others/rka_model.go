package others

import (
	"mime/multipart"
	"time"
)

type RKAModel struct {
	ID              int64      `json:"id" xml:"id"`
	Nomor           string     `json:"nomor" xml:"nomor"`
	InstansiPemberi string     `json:"instansi_pemberi" xml:"instansi_pemberi"`
	Tanggal         time.Time  `json:"tanggal" xml:"tanggal"`
	Klasifikasi     string     `json:"klasifikasi" xml:"klasifikasi"`
	MasaBerlaku     *time.Time `json:"masa_berlaku" xml:"masa_berlaku"`
	File            string     `json:"file" xml:"file"`
	Kualifikasi     int32      `json:"kualifikasi" xml:"kualifikasi"`
	IDBumd          int32      `json:"id_bumd" xml:"id_bumd"`
	IsSeumurHidup   int32      `json:"is_seumur_hidup" xml:"is_seumur_hidup"`
}

type RKAForm struct {
	Nomor           string                `json:"nomor" xml:"nomor" form:"nomor"`
	InstansiPemberi string                `json:"instansi_pemberi" xml:"instansi_pemberi" form:"instansi_pemberi"`
	Tanggal         string                `json:"tanggal" xml:"tanggal" form:"tanggal"`
	Klasifikasi     string                `json:"klasifikasi" xml:"klasifikasi" form:"klasifikasi"`
	MasaBerlaku     *string               `json:"masa_berlaku" xml:"masa_berlaku" form:"masa_berlaku"`
	File            *multipart.FileHeader `json:"file" xml:"file" form:"file"`
	Kualifikasi     int32                 `json:"kualifikasi" xml:"kualifikasi" form:"kualifikasi"`
	IDBumd          int32                 `json:"id_bumd" xml:"id_bumd" form:"id_bumd"`
}
