package others

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type ProdukModel struct {
	Id         uuid.UUID `json:"id"`
	IdBumd     uuid.UUID `json:"id_bumd"`
	NamaProduk string    `json:"nama_produk"`
	Deskripsi  string    `json:"deskripsi"`
	FotoProduk string    `json:"foto_produk"`
}

type ProdukForm struct {
	NamaProduk string                `json:"nama_produk" form:"nama_produk"`
	Deskripsi  string                `json:"deskripsi" form:"deskripsi"`
	FotoProduk *multipart.FileHeader `json:"foto_produk" form:"foto_produk"`
}
