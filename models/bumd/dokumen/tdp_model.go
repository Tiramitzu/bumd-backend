package dokumen

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type TdpModel struct {
	Id              uuid.UUID  `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd          uuid.UUID  `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nomor           string     `json:"nomor" xml:"nomor" example:"1234567890"`
	InstansiPemberi string     `json:"instansi_pemberi" xml:"instansi_pemberi" example:"Kemendagri"`
	Tanggal         time.Time  `json:"tanggal" xml:"tanggal"`
	Klasifikasi     string     `json:"klasifikasi" xml:"klasifikasi" example:"TDP"`
	MasaBerlaku     *time.Time `json:"masa_berlaku" xml:"masa_berlaku" example:"2021-01-01T00:00:00Z"`
	File            string     `json:"file" xml:"file" example:"/path/to/file.pdf"`
	Kualifikasi     int32      `json:"kualifikasi" xml:"kualifikasi" example:"1"`
	IsSeumurHidup   int32      `json:"is_seumur_hidup" xml:"is_seumur_hidup" example:"1"`
}

type TdpForm struct {
	Nomor           string                `json:"nomor" xml:"nomor" form:"nomor" example:"1234567890"`
	InstansiPemberi string                `json:"instansi_pemberi" xml:"instansi_pemberi" form:"instansi_pemberi" example:"Kemendagri"`
	Tanggal         string                `json:"tanggal" xml:"tanggal" form:"tanggal" example:"2021-01-01T00:00:00Z"`
	Klasifikasi     string                `json:"klasifikasi" xml:"klasifikasi" form:"klasifikasi" example:"TDP"`
	MasaBerlaku     *string               `json:"masa_berlaku" xml:"masa_berlaku" form:"masa_berlaku" example:"2021-01-01T00:00:00Z"`
	File            *multipart.FileHeader `json:"file" xml:"file" form:"file"`
	Kualifikasi     int32                 `json:"kualifikasi" xml:"kualifikasi" form:"kualifikasi" validate:"required" example:"1"`
}
