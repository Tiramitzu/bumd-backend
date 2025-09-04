package models

type PendidikanModel struct {
	ID              int64  `json:"id" xml:"id"`
	IDBumd          int32  `json:"id_bumd" xml:"id_bumd"`
	JumlahPegawai   int32  `json:"jumlah_pegawai" xml:"jumlah_pegawai"`
	Nomor           string `json:"nomor" xml:"nomor"`
	InstansiPemberi string `json:"instansi_pemberi" xml:"instansi_pemberi"`
	Tanggal         string `json:"tanggal" xml:"tanggal"`
	Kualifikasi     string `json:"kualifikasi" xml:"kualifikasi"`
	Klasifikasi     string `json:"klasifikasi" xml:"klasifikasi"`
	MasaBerlaku     string `json:"masa_berlaku" xml:"masa_berlaku"`
	File            string `json:"file" xml:"file"`
}

type PendidikanForm struct {
	IDBumd          int32  `json:"id_bumd" xml:"id_bumd"`
	JumlahPegawai   int32  `json:"jumlah_pegawai" xml:"jumlah_pegawai"`
	Nomor           string `json:"nomor" xml:"nomor"`
	InstansiPemberi string `json:"instansi_pemberi" xml:"instansi_pemberi"`
	Tanggal         string `json:"tanggal" xml:"tanggal"`
	Kualifikasi     string `json:"kualifikasi" xml:"kualifikasi"`
	Klasifikasi     string `json:"klasifikasi" xml:"klasifikasi"`
	MasaBerlaku     string `json:"masa_berlaku" xml:"masa_berlaku"`
	File            string `json:"file" xml:"file"`
}
