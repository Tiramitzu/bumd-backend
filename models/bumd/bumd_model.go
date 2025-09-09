package bumd

type BumdModel struct {
	ID                  int64  `json:"id" xml:"id"`
	IDDaerah            int32  `json:"id_daerah" xml:"id_daerah"`
	NamaDaerah          string `json:"nama_daerah" xml:"nama_daerah"`
	IDProvinsi          int32  `json:"id_provinsi" xml:"id_provinsi"`
	NamaProvinsi        string `json:"nama_provinsi" xml:"nama_provinsi"`
	IDBentukHukum       int32  `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha       int32  `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IDIndukPerusahaan   int32  `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	NamaIndukPerusahaan string `json:"nama_induk_perusahaan" xml:"nama_induk_perusahaan"`
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
	PenerapanSPI        bool   `json:"penerapan_spi" xml:"penerapan_spi"`
}

type BumdForm struct {
	Nama              string `json:"nama" xml:"nama"`
	Deskripsi         string `json:"deskripsi" xml:"deskripsi"`
	Alamat            string `json:"alamat" xml:"alamat"`
	NoTelp            string `json:"no_telp" xml:"no_telp"`
	NoFax             string `json:"no_fax" xml:"no_fax"`
	Email             string `json:"email" xml:"email"`
	Website           string `json:"website" xml:"website"`
	Narahubung        string `json:"narahubung" xml:"narahubung"`
	IDBentukHukum     int32  `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IDBentukUsaha     int32  `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IDIndukPerusahaan int32  `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	PenerapanSPI      bool   `json:"penerapan_spi" xml:"penerapan_spi"`
}
