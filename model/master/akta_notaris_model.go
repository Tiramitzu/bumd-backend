package models

type AktaNotarisModel struct {
	ID           int     `json:"id" xml:"id"`
	IDBumd       int     `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
	CreatedAt    string  `json:"created_at" xml:"created_at"`
	CreatedBy    int     `json:"created_by" xml:"created_by"`
	UpdatedAt    string  `json:"updated_at" xml:"updated_at"`
	UpdatedBy    int     `json:"updated_by" xml:"updated_by"`
	DeletedAt    string  `json:"deleted_at" xml:"deleted_at"`
	DeletedBy    int     `json:"deleted_by" xml:"deleted_by"`
}

type AktaNotarisForm struct {
	IDBumd       int     `json:"id_bumd" xml:"id_bumd"`
	NomorPerda   string  `json:"nomor_perda" xml:"nomor_perda"`
	TanggalPerda string  `json:"tanggal_perda" xml:"tanggal_perda"`
	ModalDasar   float64 `json:"modal_dasar" xml:"modal_dasar"`
	Keterangan   string  `json:"keterangan" xml:"keterangan"`
	File         string  `json:"file" xml:"file"`
}
