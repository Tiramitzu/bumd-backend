package keuangan

import (
	"time"

	"github.com/google/uuid"
)

type KeuModalModel struct {
	Id         uuid.UUID `json:"id"`
	IdBumd     uuid.UUID `json:"id_bumd"`
	IdProv     int64     `json:"id_prov"`
	NamaProv   string    `json:"nama_prov"`
	IdKab      int64     `json:"id_kab"`
	NamaKab    string    `json:"nama_kab"`
	Pemegang   string    `json:"pemegang"`
	NoBa       string    `json:"no_ba"`
	Tanggal    time.Time `json:"tanggal"`
	Jumlah     float64   `json:"jumlah"`
	Keterangan string    `json:"keterangan"`
}

type KeuModalForm struct {
	IdProv     int64   `json:"id_prov"`
	IdKab      int64   `json:"id_kab"`
	Pemegang   string  `json:"pemegang"`
	NoBa       string  `json:"no_ba"`
	Tanggal    string  `json:"tanggal"`
	Jumlah     float64 `json:"jumlah"`
	Keterangan string  `json:"keterangan"`
}
