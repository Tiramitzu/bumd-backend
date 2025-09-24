package models

import (
	"github.com/google/uuid"
)

type BisnisMatchingModel struct {
	IdBumd     uuid.UUID     `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	NamaBumd   string        `json:"nama_bumd" xml:"nama_bumd" example:"BUMD 1"`
	LogoBumd   string        `json:"logo_bumd" xml:"logo_bumd" example:"/path/to/file.png"`
	ProdukShow []ProdukModel `json:"produk_show" xml:"produk_show"`
	Produk     []ProdukModel `json:"produk" xml:"produk"`
}

type BumdModel struct {
	IdBumd   uuid.UUID `json:"id_bumd" xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	NamaBumd string    `json:"nama_bumd" xml:"nama_bumd" example:"BUMD 1"`
	LogoBumd string    `json:"logo_bumd" xml:"logo_bumd" example:"/path/to/file.png"`
}

type ProdukModel struct {
	IdProduk   uuid.UUID `json:"id_produk" xml:"id_produk" example:"123e4567-e89b-12d3-a456-426614174000"`
	NamaProduk string    `json:"nama_produk" xml:"nama_produk" example:"Produk 1"`
	FotoProduk string    `json:"foto_produk" xml:"foto_produk" example:"/path/to/file.png"`
	IsShow     int       `json:"is_show" xml:"is_show" example:"0"`
}
