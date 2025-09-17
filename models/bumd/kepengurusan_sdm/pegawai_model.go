package kepengurusan_sdm

import "github.com/google/uuid"

type PegawaiModel struct {
	Id            uuid.UUID `json:"id"`
	IdBumd        uuid.UUID `json:"id_bumd"`
	Tahun         int       `json:"tahun"`
	StatusPegawai int       `json:"status_pegawai"`
	Pendidikan    uuid.UUID `json:"pendidikan"`
	JumlahPegawai int       `json:"jumlah_pegawai"`
}

type PegawaiForm struct {
	Tahun         int       `json:"tahun" form:"tahun"`
	StatusPegawai int       `json:"status_pegawai" form:"status_pegawai"`
	Pendidikan    uuid.UUID `json:"pendidikan" form:"pendidikan"`
	JumlahPegawai int       `json:"jumlah_pegawai" form:"jumlah_pegawai"`
}
