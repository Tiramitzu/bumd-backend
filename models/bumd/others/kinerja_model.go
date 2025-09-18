package others

import (
	"time"

	"github.com/google/uuid"
)

type KinerjaModel struct {
	IdKinerja           uuid.UUID `json:"id_kinerja" example:"123e4567-e89b-12d3-a456-426614174000"`
	IdBumd              uuid.UUID `json:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tahun               int64     `json:"tahun" example:"2021"`
	Ebit                float64   `json:"ebit" example:"1000000"`
	Ebitda              float64   `json:"ebitda" example:"1000000"`
	ModalSendiri        float64   `json:"modal_sendiri" example:"1000000"`
	Penyusutan          float64   `json:"penyusutan" example:"1000000"`
	CapitalEmployed     float64   `json:"capital_employed" example:"1000000"`
	TotalAsetAwal       float64   `json:"total_aset_awal" example:"1000000"`
	TotalAsetAkhir      float64   `json:"total_aset_akhir" example:"1000000"`
	Kas                 float64   `json:"kas" example:"1000000"`
	SetaraKas           float64   `json:"setara_kas" example:"1000000"`
	KewajibanLancar     float64   `json:"kewajiban_lancar" example:"1000000"`
	Hartalancar         float64   `json:"harta_lancar" example:"1000000"`
	PenjualanBersih     float64   `json:"penjualan_bersih" example:"1000000"`
	PiutangDagang       float64   `json:"piutang_dagang" example:"1000000"`
	HargaPokokPenjualan float64   `json:"harga_pokok_penjualan" example:"1000000"`
	Persediaan          float64   `json:"persediaan" example:"1000000"`
	AktivaTetap         float64   `json:"aktiva_tetap" example:"1000000"`
	AkumulasiDepresiasi float64   `json:"akumulasi_depresiasi" example:"1000000"`
	KreditBermasalah    float64   `json:"kredit_bermasalah" example:"1000000"`
	TotalKredit         float64   `json:"total_kredit" example:"1000000"`
	CreatedAt           time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy           int64     `json:"created_by" example:"1"`
	UpdatedAt           time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy           int64     `json:"updated_by" example:"1"`
}

type KinerjaForm struct {
	Tahun               int64   `json:"tahun" validate:"required,min=2000,max=2099" example:"2021"`
	Ebit                float64 `json:"ebit" validate:"required,min=0" example:"1000000"`
	Ebitda              float64 `json:"ebitda" validate:"required,min=0" example:"1000000"`
	ModalSendiri        float64 `json:"modal_sendiri" validate:"required,min=0" example:"1000000"`
	Penyusutan          float64 `json:"penyusutan" validate:"required,min=0" example:"1000000"`
	CapitalEmployed     float64 `json:"capital_employed" validate:"required,min=0" example:"1000000"`
	TotalAsetAwal       float64 `json:"total_aset_awal" validate:"required,min=0" example:"1000000"`
	TotalAsetAkhir      float64 `json:"total_aset_akhir" validate:"required,min=0" example:"1000000"`
	Kas                 float64 `json:"kas" validate:"required,min=0" example:"1000000"`
	SetaraKas           float64 `json:"setara_kas" validate:"required,min=0" example:"1000000"`
	KewajibanLancar     float64 `json:"kewajiban_lancar" validate:"required,min=0" example:"1000000"`
	Hartalancar         float64 `json:"harta_lancar" validate:"required,min=0" example:"1000000"`
	PenjualanBersih     float64 `json:"penjualan_bersih" validate:"required,min=0" example:"1000000"`
	PiutangDagang       float64 `json:"piutang_dagang" validate:"required,min=0" example:"1000000"`
	HargaPokokPenjualan float64 `json:"harga_pokok_penjualan" validate:"required,min=0" example:"1000000"`
	Persediaan          float64 `json:"persediaan" validate:"required,min=0" example:"1000000"`
	AktivaTetap         float64 `json:"aktiva_tetap" validate:"required,min=0" example:"1000000"`
	AkumulasiDepresiasi float64 `json:"akumulasi_depresiasi" validate:"required,min=0" example:"1000000"`
	KreditBermasalah    float64 `json:"kredit_bermasalah" validate:"required,min=0" example:"1000000"`
	TotalKredit         float64 `json:"total_kredit" validate:"required,min=0" example:"1000000"`
}
