package models

type ProdukModel struct {
	ID        int    `json:"id" xml:"id"`
	IDBumd    int    `json:"id_bumd" xml:"id_bumd"`
	IDDaerah  int    `json:"id_daerah" xml:"id_daerah"`
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
	Gambar    string `json:"gambar" xml:"gambar"`
}

type ProdukForm struct {
	IDBumd    int    `json:"id_bumd" xml:"id_bumd"`
	IDDaerah  int    `json:"id_daerah" xml:"id_daerah"`
	Nama      string `json:"nama" xml:"nama"`
	Deskripsi string `json:"deskripsi" xml:"deskripsi"`
	Gambar    string `json:"gambar" xml:"gambar"`
}
