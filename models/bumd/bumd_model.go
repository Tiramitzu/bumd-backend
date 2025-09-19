package bumd

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type BumdModel struct {
	Id                  uuid.UUID  `json:"id" xml:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBentukHukum       uuid.UUID  `json:"id_bentuk_hukum" xml:"id_bentuk_hukum" example:"01994c00-6768-7fde-ab34-f315eea6510f"`
	IdBentukUsaha       uuid.UUID  `json:"id_bentuk_usaha" xml:"id_bentuk_usaha" example:"01994c01-c285-7e57-a486-fd9978083917"`
	IdIndukPerusahaan   *uuid.UUID `json:"id_induk_perusahaan" xml:"id_induk_perusahaan" example:"123e4567-e89b-12d3-a456-426614174000"`
	NamaDaerah          string     `json:"nama_daerah" xml:"nama_daerah" example:"DKI Jakarta"`
	NamaIndukPerusahaan string     `json:"nama_induk_perusahaan" xml:"nama_induk_perusahaan" example:"PT. BUMD"`
	BentukBadanHukum    string     `json:"bentuk_badan_hukum" xml:"bentuk_badan_hukum" example:"Perumda"`
	BentukUsaha         string     `json:"bentuk_usaha" xml:"bentuk_usaha" example:"Air Minum"`
	Nama                string     `json:"nama" xml:"nama" example:"BUMD"`
	Deskripsi           string     `json:"deskripsi" xml:"deskripsi" example:"BUMD adalah Badan Usaha Milik Daerah"`
	Alamat              string     `json:"alamat" xml:"alamat" example:"Jl. Raya No. 1, Jakarta"`
	NoTelp              string     `json:"no_telp" xml:"no_telp" example:"081234567890"`
	NoFax               string     `json:"no_fax" xml:"no_fax" example:"081234567890"`
	Email               string     `json:"email" xml:"email" example:"bumd@e-bumd.com"`
	Website             string     `json:"website" xml:"website" example:"https://e-bumd.com"`
	Narahubung          string     `json:"narahubung" xml:"narahubung" example:"John Doe"`
	NPWP                string     `json:"npwp" xml:"npwp" example:"123456789012345"`
	NPWPPemberi         string     `json:"npwp_pemberi" xml:"npwp_pemberi" example:"DJP"`
	NPWPFile            string     `json:"npwp_file" xml:"npwp_file" example:"/path/to/file.pdf"`
	SPIFile             string     `json:"file_spi" xml:"file_spi" example:"/path/to/file.pdf"`
	Logo                string     `json:"logo" xml:"logo" example:"/path/to/file.png"`
	IdDaerah            int32      `json:"id_daerah" xml:"id_daerah" example:"1"`
	IdProvinsi          int32      `json:"id_provinsi" xml:"id_provinsi" example:"1"`
	NamaProvinsi        string     `json:"nama_provinsi" xml:"nama_provinsi" example:"DKI Jakarta"`
	PenerapanSPI        bool       `json:"penerapan_spi" xml:"penerapan_spi" example:"true"`
	CreatedAt           time.Time  `json:"created_at" xml:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy           int64      `json:"created_by" xml:"created_by" example:"1"`
	UpdatedAt           time.Time  `json:"updated_at" xml:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy           int64      `json:"updated_by" xml:"updated_by" example:"1"`
}

type BumdForm struct {
	IdBentukHukum     string `json:"id_bentuk_hukum" xml:"id_bentuk_hukum" validate:"required,custom_uuid" example:"01994c00-6768-7fde-ab34-f315eea6510f"`
	IdBentukUsaha     string `json:"id_bentuk_usaha" xml:"id_bentuk_usaha" validate:"required,custom_uuid" example:"01994c01-c285-7e57-a486-fd9978083917"`
	IdIndukPerusahaan string `json:"id_induk_perusahaan" xml:"id_induk_perusahaan" validate:"required,custom_uuid" example:"00000000-0000-0000-0000-000000000000" default:"00000000-0000-0000-0000-000000000000"`
	Nama              string `json:"nama" xml:"nama" validate:"required" example:"BUMD"`
	Deskripsi         string `json:"deskripsi" xml:"deskripsi" example:"BUMD adalah Badan Usaha Milik Daerah"`
	Alamat            string `json:"alamat" xml:"alamat" example:"Jl. Raya No. 1, Jakarta"`
	NoTelp            string `json:"no_telp" xml:"no_telp" example:"081234567890"`
	NoFax             string `json:"no_fax" xml:"no_fax" example:"081234567890"`
	Email             string `json:"email" xml:"email" validate:"email" example:"bumd@e-bumd.com"`
	Website           string `json:"website" xml:"website" example:"https://e-bumd.com"`
	Narahubung        string `json:"narahubung" xml:"narahubung" example:"John Doe"`
	PenerapanSPI      bool   `json:"penerapan_spi" xml:"penerapan_spi" validate:"boolean" example:"false"`
}

type LogoModel struct {
	Logo string `json:"logo" xml:"logo" example:"/path/to/file.png"`
}

type LogoForm struct {
	Logo *multipart.FileHeader `json:"logo" xml:"logo" form:"logo" validate:"required"`
}

type SPIModel struct {
	PenerapanSPI bool   `json:"penerapan_spi" xml:"penerapan_spi" example:"true"`
	FileSPI      string `json:"file_spi" xml:"file_spi" example:"/path/to/file.pdf"`
}

type SPIForm struct {
	PenerapanSPI bool                  `json:"penerapan_spi" xml:"penerapan_spi" form:"penerapan_spi" validate:"required,boolean" example:"true"`
	FileSPI      *multipart.FileHeader `json:"file_spi" xml:"file_spi" form:"file_spi"`
}

type NPWPModel struct {
	NPWP    string `json:"npwp" xml:"npwp" example:"123456789012345"`
	Pemberi string `json:"pemberi" xml:"pemberi" example:"DJP"`
	File    string `json:"file" xml:"file" example:"/path/to/file.pdf"`
}

type NPWPForm struct {
	NPWP    string                `json:"npwp" xml:"npwp" form:"npwp" validate:"required" example:"123456789012345"`
	Pemberi string                `json:"pemberi" xml:"pemberi" form:"pemberi" validate:"required" example:"DJP"`
	File    *multipart.FileHeader `json:"file" xml:"file" form:"file"`
}
