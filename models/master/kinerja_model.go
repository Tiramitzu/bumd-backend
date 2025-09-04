package models

type KinerjaModel struct {
	ID                       int64   `json:"id" xml:"id"`
	IDBumd                   int32   `json:"id_bumd" xml:"id_bumd"`
	Tahun                    int32   `json:"tahun" xml:"tahun"`
	LabaBersihSebelumPajak   float64 `json:"laba_bersih_sebelum_pajak" xml:"laba_bersih_sebelum_pajak"`
	LabaBersihSetelahPajak   float64 `json:"laba_bersih_setelah_pajak" xml:"laba_bersih_setelah_pajak"`
	ModalSendiri             float64 `json:"modal_sendiri" xml:"modal_sendiri"`
	Penyusutan               float64 `json:"penyusutan" xml:"penyusutan"`
	TotalAsetKewajibanLancar float64 `json:"total_aset_kewajiban_lancar" xml:"total_aset_kewajiban_lancar"`
	TotalAsetAwal            float64 `json:"total_aset_awal" xml:"total_aset_awal"`
	TotalAsetAkhir           float64 `json:"total_aset_akhir" xml:"total_aset_akhir"`
	Kas                      float64 `json:"kas" xml:"kas"`
	SetaraKas                float64 `json:"setara_kas" xml:"setara_kas"`
	KewajibanLancar          float64 `json:"kewajiban_lancar" xml:"kewajiban_lancar"`
	HartaLancar              float64 `json:"harta_lancar" xml:"harta_lancar"`
	PenjualanBersih          float64 `json:"penjualan_bersih" xml:"penjualan_bersih"`
	RataRataPiutangDagang    float64 `json:"rata_rata_piutang_dagang" xml:"rata_rata_piutang_dagang"`
	HargaPokokPenjualan      float64 `json:"harga_pokok_penjualan" xml:"harga_pokok_penjualan"`
	RataRataPersediaan       float64 `json:"rata_rata_persediaan" xml:"rata_rata_persediaan"`
	AktivaTetap              float64 `json:"aktiva_tetap" xml:"aktiva_tetap"`
	AkumulasiDepresiasi      float64 `json:"akumulasi_depresiasi" xml:"akumulasi_depresiasi"`
	KreditBermasalah         float64 `json:"kredit_bermasalah" xml:"kredit_bermasalah"`
	TotalKredit              float64 `json:"total_kredit" xml:"total_kredit"`
}

type KinerjaForm struct {
	IDBumd                   int32   `json:"id_bumd" xml:"id_bumd"`
	Tahun                    int32   `json:"tahun" xml:"tahun"`
	LabaBersihSebelumPajak   float64 `json:"laba_bersih_sebelum_pajak" xml:"laba_bersih_sebelum_pajak"`
	LabaBersihSetelahPajak   float64 `json:"laba_bersih_setelah_pajak" xml:"laba_bersih_setelah_pajak"`
	ModalSendiri             float64 `json:"modal_sendiri" xml:"modal_sendiri"`
	Penyusutan               float64 `json:"penyusutan" xml:"penyusutan"`
	TotalAsetKewajibanLancar float64 `json:"total_aset_kewajiban_lancar" xml:"total_aset_kewajiban_lancar"`
	TotalAsetAwal            float64 `json:"total_aset_awal" xml:"total_aset_awal"`
	TotalAsetAkhir           float64 `json:"total_aset_akhir" xml:"total_aset_akhir"`
	Kas                      float64 `json:"kas" xml:"kas"`
	SetaraKas                float64 `json:"setara_kas" xml:"setara_kas"`
	KewajibanLancar          float64 `json:"kewajiban_lancar" xml:"kewajiban_lancar"`
	HartaLancar              float64 `json:"harta_lancar" xml:"harta_lancar"`
	PenjualanBersih          float64 `json:"penjualan_bersih" xml:"penjualan_bersih"`
	RataRataPiutangDagang    float64 `json:"rata_rata_piutang_dagang" xml:"rata_rata_piutang_dagang"`
	HargaPokokPenjualan      float64 `json:"harga_pokok_penjualan" xml:"harga_pokok_penjualan"`
	RataRataPersediaan       float64 `json:"rata_rata_persediaan" xml:"rata_rata_persediaan"`
	AktivaTetap              float64 `json:"aktiva_tetap" xml:"aktiva_tetap"`
	AkumulasiDepresiasi      float64 `json:"akumulasi_depresiasi" xml:"akumulasi_depresiasi"`
	KreditBermasalah         float64 `json:"kredit_bermasalah" xml:"kredit_bermasalah"`
	TotalKredit              float64 `json:"total_kredit" xml:"total_kredit"`
}
