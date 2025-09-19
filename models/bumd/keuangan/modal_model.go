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
	CreatedAt  time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy  int64     `json:"created_by" example:"1"`
	UpdatedAt  time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy  int64     `json:"updated_by" example:"1"`
}

type KeuModalForm struct {
	IdProv     int64   `json:"id_prov" validate:"min=0" example:"0"`
	IdKab      int64   `json:"id_kab" validate:"min=0" example:"0"`
	Pemegang   string  `json:"pemegang" validate:"required_without_all=IdProv IdKab" example:"Pemegang Modal"`
	NoBa       string  `json:"no_ba" validate:"required" example:"1234567890"`
	Tanggal    string  `json:"tanggal" validate:"required,datetime=2006-01-02" example:"2021-01-01"`
	Jumlah     float64 `json:"jumlah" validate:"min=0" example:"1000000"`
	Keterangan string  `json:"keterangan" example:"Keterangan Modal"`
}
