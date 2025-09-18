package keuangan

import (
	"time"

	"github.com/google/uuid"
)

type KeuModalModel struct {
	Id         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd     uuid.UUID `json:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdProv     int64     `json:"id_prov" example:"1"`
	NamaProv   string    `json:"nama_prov" example:"DKI Jakarta"`
	IdKab      int64     `json:"id_kab" example:"1"`
	NamaKab    string    `json:"nama_kab" example:"Jakarta Selatan"`
	Pemegang   string    `json:"pemegang" example:"Pemegang Modal"`
	NoBa       string    `json:"no_ba" example:"1234567890"`
	Tanggal    time.Time `json:"tanggal" example:"2021-01-01T00:00:00Z"`
	Jumlah     float64   `json:"jumlah" example:"1000000"`
	Keterangan string    `json:"keterangan" example:"Keterangan Modal"`
}

type KeuModalForm struct {
	IdProv     int64   `json:"id_prov" validate:"required" example:"1"`
	IdKab      int64   `json:"id_kab" validate:"required" example:"1"`
	Pemegang   string  `json:"pemegang" example:"Pemegang Modal"`
	NoBa       string  `json:"no_ba" example:"1234567890"`
	Tanggal    string  `json:"tanggal" validate:"required,datetime=2006-01-02" example:"2021-01-01T00:00:00Z"`
	Jumlah     float64 `json:"jumlah" validate:"required,min=0" example:"1000000"`
	Keterangan string  `json:"keterangan" example:"Keterangan Modal"`
}
