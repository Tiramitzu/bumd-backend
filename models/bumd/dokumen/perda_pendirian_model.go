package dokumen

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PerdaPendirianModel struct {
	Id         uuid.UUID `json:"id" xml:"id"`
	IdBumd     uuid.UUID `json:"id_bumd" xml:"id_bumd"`
	Nomor      string    `json:"nomor" xml:"nomor"`
	Tanggal    time.Time `json:"tanggal" xml:"tanggal"`
	Keterangan string    `json:"keterangan" xml:"keterangan"`
	File       string    `json:"file" xml:"file"`
	ModalDasar float64   `json:"modal_dasar" xml:"modal_dasar"`
}

type PerdaPendirianForm struct {
	Nomor      string                `json:"nomor_perda" xml:"nomor_perda" form:"nomor_perda"`
	Tanggal    string                `json:"tanggal_perda" xml:"tanggal_perda" form:"tanggal_perda"`
	Keterangan string                `json:"keterangan" xml:"keterangan" form:"keterangan"`
	File       *multipart.FileHeader `json:"file" xml:"file" form:"file"`
	ModalDasar string                `json:"modal_dasar" xml:"modal_dasar" form:"modal_dasar"`
}
