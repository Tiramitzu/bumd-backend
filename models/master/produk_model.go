package models

type ProdukModel struct {
	ID        int64   `json:"id" xml:"id"`
	IDBumd    int32   `json:"id_bumd" xml:"id_bumd"`
	IDDaerah  int32   `json:"id_daerah" xml:"id_daerah"`
	Nama      string  `json:"nama" xml:"nama"`
	Deskripsi string  `json:"deskripsi" xml:"deskripsi"`
	Nilai     float64 `json:"nilai" xml:"nilai"`
	Gambar    string  `json:"gambar" xml:"gambar"`
	Tampilkan bool    `json:"tampilkan" xml:"tampilkan"`
}

type ProdukForm struct {
	IDBumd    int32   `json:"id_bumd" xml:"id_bumd"`
	IDDaerah  int32   `json:"id_daerah" xml:"id_daerah"`
	Nama      string  `json:"nama" xml:"nama"`
	Deskripsi string  `json:"deskripsi" xml:"deskripsi"`
	Nilai     float64 `json:"nilai" xml:"nilai"`
	Gambar    string  `json:"gambar" xml:"gambar"`
	Tampilkan bool    `json:"tampilkan" xml:"tampilkan"`
}
