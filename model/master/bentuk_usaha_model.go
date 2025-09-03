package models

type BentukUsahaModel struct {
	ID        int    `json:"id" xml:"id"`
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}

type BentukUsahaForm struct {
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}
