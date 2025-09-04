package models

type BentukBadanHukumModel struct {
	ID        int64  `json:"id" xml:"id"`
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}

type BentukBadanHukumForm struct {
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
}
