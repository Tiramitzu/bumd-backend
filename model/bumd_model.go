package models

type BumdModel struct {
	ID                  int    `json:"id" xml:"id"`
	IDDaerah            int    `json:"id_daerah" xml:"id_daerah"`
	IDBentukHukum       int    `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha       int    `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IDIndukPerusahaan   int    `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	NamaIndukPerusahaan string `json:"nama_induk_perusahaan" xml:"nama_induk_perusahaan"`
	PenerapanSPI        bool   `json:"penerapan_spi" xml:"penerapan_spi"`
	BentukBadanHukum    string `json:"bentuk_badan_hukum" xml:"bentuk_badan_hukum"`
	BentukUsaha         string `json:"bentuk_usaha" xml:"bentuk_usaha"`
	Nama                string `json:"nama" xml:"nama"`
	Deskripsi           string `json:"deskripsi" xml:"deskripsi"`
	Alamat              string `json:"alamat" xml:"alamat"`
	NoTelp              string `json:"no_telp" xml:"no_telp"`
	NoFax               string `json:"no_fax" xml:"no_fax"`
	Email               string `json:"email" xml:"email"`
	Website             string `json:"website" xml:"website"`
	Narahubung          string `json:"narahubung" xml:"narahubung"`
}

type BumdForm struct {
	IDBentukHukum     int    `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha     int    `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IDIndukPerusahaan int    `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	PenerapanSPI      bool   `json:"penerapan_spi" xml:"penerapan_spi"`
	Nama              string `json:"nama" xml:"nama"`
	Deskripsi         string `json:"deskripsi" xml:"deskripsi"`
	Alamat            string `json:"alamat" xml:"alamat"`
	NoTelp            string `json:"no_telp" xml:"no_telp"`
	NoFax             string `json:"no_fax" xml:"no_fax"`
	Email             string `json:"email" xml:"email"`
	Website           string `json:"website" xml:"website"`
	Narahubung        string `json:"narahubung" xml:"narahubung"`
}
