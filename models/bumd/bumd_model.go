package bumd

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type BumdModel struct {
	Id                  uuid.UUID  `json:"id" xml:"id"`
	IdBentukHukum       uuid.UUID  `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IdBentukUsaha       uuid.UUID  `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IdIndukPerusahaan   *uuid.UUID `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	NamaDaerah          string     `json:"nama_daerah" xml:"nama_daerah"`
	NamaIndukPerusahaan string     `json:"nama_induk_perusahaan" xml:"nama_induk_perusahaan"`
	BentukBadanHukum    string     `json:"bentuk_badan_hukum" xml:"bentuk_badan_hukum"`
	BentukUsaha         string     `json:"bentuk_usaha" xml:"bentuk_usaha"`
	Nama                string     `json:"nama" xml:"nama"`
	Deskripsi           string     `json:"deskripsi" xml:"deskripsi"`
	Alamat              string     `json:"alamat" xml:"alamat"`
	NoTelp              string     `json:"no_telp" xml:"no_telp"`
	NoFax               string     `json:"no_fax" xml:"no_fax"`
	Email               string     `json:"email" xml:"email"`
	Website             string     `json:"website" xml:"website"`
	Narahubung          string     `json:"narahubung" xml:"narahubung"`
	NPWP                string     `json:"npwp" xml:"npwp"`
	NPWPPemberi         string     `json:"npwp_pemberi" xml:"npwp_pemberi"`
	NPWPFile            string     `json:"npwp_file" xml:"npwp_file"`
	SPIFile             string     `json:"file_spi" xml:"file_spi"`
	Logo                string     `json:"logo" xml:"logo"`
	IdDaerah            int32      `json:"id_daerah" xml:"id_daerah"`
	IdProvinsi          int32      `json:"id_provinsi" xml:"id_provinsi"`
	NamaProvinsi        string     `json:"nama_provinsi" xml:"nama_provinsi"`
	PenerapanSPI        bool       `json:"penerapan_spi" xml:"penerapan_spi"`
	CreatedAt           time.Time  `json:"created_at" xml:"created_at"`
	CreatedBy           int64      `json:"created_by" xml:"created_by"`
	UpdatedAt           time.Time  `json:"updated_at" xml:"updated_at"`
	UpdatedBy           int64      `json:"updated_by" xml:"updated_by"`
}

type BumdForm struct {
	IdBentukHukum     uuid.UUID `json:"id_bentuk_hukum" xml:"id_bentuk_hukum"`
	IdBentukUsaha     uuid.UUID `json:"id_bentuk_usaha" xml:"id_bentuk_usaha"`
	IdIndukPerusahaan uuid.UUID `json:"id_induk_perusahaan" xml:"id_induk_perusahaan"`
	Nama              string    `json:"nama" xml:"nama"`
	Deskripsi         string    `json:"deskripsi" xml:"deskripsi"`
	Alamat            string    `json:"alamat" xml:"alamat"`
	NoTelp            string    `json:"no_telp" xml:"no_telp"`
	NoFax             string    `json:"no_fax" xml:"no_fax"`
	Email             string    `json:"email" xml:"email"`
	Website           string    `json:"website" xml:"website"`
	Narahubung        string    `json:"narahubung" xml:"narahubung"`
	PenerapanSPI      bool      `json:"penerapan_spi" xml:"penerapan_spi"`
}

type SPIModel struct {
	PenerapanSPI bool   `json:"penerapan_spi" xml:"penerapan_spi"`
	FileSPI      string `json:"file_spi" xml:"file_spi"`
}

type SPIForm struct {
	PenerapanSPI bool                  `json:"penerapan_spi" xml:"penerapan_spi"`
	FileSPI      *multipart.FileHeader `json:"file_spi" xml:"file_spi"`
}

type NPWPModel struct {
	NPWP    string `json:"npwp" xml:"npwp"`
	Pemberi string `json:"pemberi" xml:"pemberi"`
	File    string `json:"file" xml:"file"`
}

type NPWPForm struct {
	NPWP    string                `json:"npwp" xml:"npwp"`
	Pemberi string                `json:"pemberi" xml:"pemberi"`
	File    *multipart.FileHeader `json:"file" xml:"file"`
}
