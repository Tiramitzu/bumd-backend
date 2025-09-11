package others

type PegawaiModel struct {
	ID            int `json:"id"`
	IDBumd        int `json:"id_bumd"`
	Tahun         int `json:"tahun"`
	StatusPegawai int `json:"status_pegawai"`
	Pendidikan    int `json:"pendidikan"`
	JumlahPegawai int `json:"jumlah_pegawai"`
}

type PegawaiForm struct {
	Tahun         int `json:"tahun" form:"tahun"`
	StatusPegawai int `json:"status_pegawai" form:"status_pegawai"`
	Pendidikan    int `json:"pendidikan" form:"pendidikan"`
	JumlahPegawai int `json:"jumlah_pegawai" form:"jumlah_pegawai"`
}
