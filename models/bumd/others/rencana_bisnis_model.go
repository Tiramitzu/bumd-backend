package others

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type RencanaBisnisModel struct {
	Id              uuid.UUID  `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd          uuid.UUID  `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nomor           string     `json:"nomor" xml:"nomor" example:"1234567890"`
	InstansiPemberi string     `json:"instansi_pemberi" xml:"instansi_pemberi" example:"Kemendagri"`
	Tanggal         time.Time  `json:"tanggal" xml:"tanggal" example:"2021-01-01T00:00:00Z"`
	Klasifikasi     string     `json:"klasifikasi" xml:"klasifikasi" example:"Rencana Bisnis"`
	MasaBerlaku     *time.Time `json:"masa_berlaku" xml:"masa_berlaku" example:"2021-01-01T00:00:00Z"`
	File            string     `json:"file" xml:"file" example:"/path/to/file.pdf"`
	Kualifikasi     int32      `json:"kualifikasi" xml:"kualifikasi" example:"1"`
	IsSeumurHidup   int32      `json:"is_seumur_hidup" xml:"is_seumur_hidup" example:"1"`
	CreatedAt       time.Time  `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy       int64      `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt       time.Time  `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy       int64      `json:"updated_by" xml:"updated_by" example:"1"`
}

type RencanaBisnisForm struct {
	Nomor           string                `json:"nomor" xml:"nomor" form:"nomor" validate:"required" example:"1234567890"`
	InstansiPemberi string                `json:"instansi_pemberi" xml:"instansi_pemberi" form:"instansi_pemberi" validate:"required" example:"Kemendagri"`
	Tanggal         string                `json:"tanggal" xml:"tanggal" form:"tanggal" validate:"required,datetime=2006-01-02" example:"2021-01-01"`
	Klasifikasi     string                `json:"klasifikasi" xml:"klasifikasi" form:"klasifikasi" example:"Rencana Bisnis"`
	MasaBerlaku     *string               `json:"masa_berlaku" xml:"masa_berlaku" form:"masa_berlaku" example:"2021-01-01T00:00:00Z"`
	File            *multipart.FileHeader `json:"file" xml:"file" form:"file"`
	Kualifikasi     int32                 `json:"kualifikasi" xml:"kualifikasi" form:"kualifikasi" validate:"oneof=0 1" example:"1"`
}
