package others

import "time"

type KinerjaModel struct {
	IdKinerja           int64     `json:"id_kinerja"`
	IdBumd              int64     `json:"id_bumd"`
	Tahun               int64     `json:"tahun"`
	Ebit                float64   `json:"ebit"`
	Ebitda              float64   `json:"ebitda"`
	ModalSendiri        float64   `json:"modal_sendiri"`
	Penyusutan          float64   `json:"penyusutan"`
	CapitalEmployed     float64   `json:"capital_employed"`
	TotalAsetAwal       float64   `json:"total_aset_awal"`
	TotalAsetAkhir      float64   `json:"total_aset_akhir"`
	Kas                 float64   `json:"kas"`
	SetaraKas           float64   `json:"setara_kas"`
	KewajibanLancar     float64   `json:"kewajiban_lancar"`
	Hartalancar         float64   `json:"harta_lancar"`
	PenjualanBersih     float64   `json:"penjualan_bersih"`
	PiutangDagang       float64   `json:"piutang_dagang"`
	HargaPokokPenjualan float64   `json:"harga_pokok_penjualan"`
	Persediaan          float64   `json:"persediaan"`
	AktivaTetap         float64   `json:"aktiva_tetap"`
	AkumulasiDepresiasi float64   `json:"akumulasi_depresiasi"`
	KreditBermasalah    float64   `json:"kredit_bermasalah"`
	TotalKredit         float64   `json:"total_kredit"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           int64     `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           int64     `json:"updated_by"`
}

type KinerjaForm struct {
	Tahun               int64   `json:"tahun"`
	Ebit                float64 `json:"ebit"`
	Ebitda              float64 `json:"ebitda"`
	ModalSendiri        float64 `json:"modal_sendiri"`
	Penyusutan          float64 `json:"penyusutan"`
	CapitalEmployed     float64 `json:"capital_employed"`
	TotalAsetAwal       float64 `json:"total_aset_awal"`
	TotalAsetAkhir      float64 `json:"total_aset_akhir"`
	Kas                 float64 `json:"kas"`
	SetaraKas           float64 `json:"setara_kas"`
	KewajibanLancar     float64 `json:"kewajiban_lancar"`
	Hartalancar         float64 `json:"harta_lancar"`
	PenjualanBersih     float64 `json:"penjualan_bersih"`
	PiutangDagang       float64 `json:"piutang_dagang"`
	HargaPokokPenjualan float64 `json:"harga_pokok_penjualan"`
	Persediaan          float64 `json:"persediaan"`
	AktivaTetap         float64 `json:"aktiva_tetap"`
	AkumulasiDepresiasi float64 `json:"akumulasi_depresiasi"`
	KreditBermasalah    float64 `json:"kredit_bermasalah"`
	TotalKredit         float64 `json:"total_kredit"`
}
