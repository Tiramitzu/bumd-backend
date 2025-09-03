package models

type BumdModel struct {
	ID               int    `json:"id" xml:"id"`
	IDDaerah         int    `json:"id_daerah" xml:"id_daerah"`
	IDBentukHukum    int    `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha    int    `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	BentukBadanHukum string `json:"bentuk_badan_hukum" xml:"bentuk_badan_hukum"`
	BentukUsaha      string `json:"bentuk_usaha" xml:"bentuk_usaha"`
	Nama             string `json:"nama" xml:"nama"`
	Deskripsi        string `json:"deskripsi" xml:"deskripsi"`
	Alamat           string `json:"alamat" xml:"alamat"`
	NoTelp           string `json:"no_telp" xml:"no_telp"`
	Email            string `json:"email" xml:"email"`
	Website          string `json:"website" xml:"website"`
}

type BumdForm struct {
	IDBentukHukum int    `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha int    `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	Nama          string `json:"nama" xml:"nama"`
	Deskripsi     string `json:"deskripsi" xml:"deskripsi"`
	Alamat        string `json:"alamat" xml:"alamat"`
	NoTelp        string `json:"no_telp" xml:"no_telp"`
	Email         string `json:"email" xml:"email"`
	Website       string `json:"website" xml:"website"`
}
