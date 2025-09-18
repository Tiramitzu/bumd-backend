package others

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type ProdukModel struct {
	Id         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd     uuid.UUID `json:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	NamaProduk string    `json:"nama_produk" example:"Produk 1"`
	Deskripsi  string    `json:"deskripsi" example:"Deskripsi Produk 1"`
	FotoProduk string    `json:"foto_produk" example:"/path/to/file.png"`
	CreatedAt  time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy  int64     `json:"created_by" example:"1"`
	UpdatedAt  time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy  int64     `json:"updated_by" example:"1"`
}

type ProdukForm struct {
	NamaProduk string                `json:"nama_produk" form:"nama_produk" validate:"required" example:"Produk 1"`
	Deskripsi  string                `json:"deskripsi" form:"deskripsi" example:"Deskripsi Produk 1"`
	FotoProduk *multipart.FileHeader `json:"foto_produk" form:"foto_produk"`
}
