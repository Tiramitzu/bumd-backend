package models

type LaporanKeuanganModel struct {
	ID           int64  `json:"id" xml:"id"`
	IDBumd       int32  `json:"id_bumd" xml:"id_bumd"`
	Tahun        int32  `json:"tahun" xml:"tahun"`
	Periode      string `json:"periode" xml:"periode"`
	JenisLaporan string `json:"jenis_laporan" xml:"jenis_laporan"`
	File         string `json:"file" xml:"file"`
	Keterangan   string `json:"keterangan" xml:"keterangan"`
}

type LaporanKeuanganForm struct {
	IDBumd       int32  `json:"id_bumd" xml:"id_bumd"`
	Tahun        int32  `json:"tahun" xml:"tahun"`
	Periode      string `json:"periode" xml:"periode"`
	JenisLaporan string `json:"jenis_laporan" xml:"jenis_laporan"`
	File         string `json:"file" xml:"file"`
	Keterangan   string `json:"keterangan" xml:"keterangan"`
}
