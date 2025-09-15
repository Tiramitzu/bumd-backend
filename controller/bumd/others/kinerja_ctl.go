package others

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type KinerjaController struct {
	pgxConn *pgxpool.Pool
}

func NewKinerjaController(pgxConn *pgxpool.Pool) *KinerjaController {
	return &KinerjaController{
		pgxConn: pgxConn,
	}
}

func (c *KinerjaController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, idBumd uuid.UUID, tahun int) (r []others.KinerjaModel, totalCount, pageCount int, err error) {
	r = make([]others.KinerjaModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, err
	}
	offset := limit * (page - 1)

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	var args []interface{}
	qCount := `
	SELECT COALESCE(COUNT(*), 0) FROM trn_kinerja WHERE deleted_by = 0 AND id_bumd = $1
	`
	args = append(args, idBumd)

	q := `
	SELECT
		trn_kinerja.id_kinerja,
		trn_kinerja.id_bumd,
		trn_kinerja.tahun_kinerja,
		trn_kinerja.ebit_kinerja,
		trn_kinerja.ebitda_kinerja,
		trn_kinerja.modal_sendiri_kinerja,
		trn_kinerja.penyusutan_kinerja,
		trn_kinerja.capital_employed_kinerja,
		trn_kinerja.total_aset_awal_kinerja,
		trn_kinerja.total_aset_akhir_kinerja,
		trn_kinerja.kas_kinerja,
		trn_kinerja.setara_kas_kinerja,
		trn_kinerja.kewajiban_lancar_kinerja,
		trn_kinerja.harta_lancar_kinerja,
		trn_kinerja.penjualan_bersih_kinerja,
		trn_kinerja.piutang_dagang_kinerja,
		trn_kinerja.harga_pokok_penjualan_kinerja,
		trn_kinerja.persediaan_kinerja,
		trn_kinerja.aktiva_tetap_kinerja,
		trn_kinerja.akumulasi_depresiasi_kinerja,
		trn_kinerja.kredit_bermasalah_kinerja,
		trn_kinerja.total_kredit_kinerja,
		trn_kinerja.created_at,
		trn_kinerja.created_by,
		trn_kinerja.updated_at,
		trn_kinerja.updated_by
	FROM trn_kinerja
	WHERE trn_kinerja.deleted_by = 0
	AND trn_kinerja.id_bumd = $1
	`
	if tahun > 0 {
		qCount += fmt.Sprintf(` AND trn_kinerja.tahun_kinerja = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND trn_kinerja.tahun_kinerja = $%d`, len(args)+1)
		args = append(args, tahun)
	}
	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, err
	}

	q += fmt.Sprintf(` ORDER BY trn_kinerja.created_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, err
	}

	for rows.Next() {
		var m others.KinerjaModel
		err = rows.Scan(
			&m.IdKinerja,
			&m.IdBumd,
			&m.Tahun,
			&m.Ebit,
			&m.Ebitda,
			&m.ModalSendiri,
			&m.Penyusutan,
			&m.CapitalEmployed,
			&m.TotalAsetAwal,
			&m.TotalAsetAkhir,
			&m.Kas,
			&m.SetaraKas,
			&m.KewajibanLancar,
			&m.Hartalancar,
			&m.PenjualanBersih,
			&m.PiutangDagang,
			&m.HargaPokokPenjualan,
			&m.Persediaan,
			&m.AktivaTetap,
			&m.AkumulasiDepresiasi,
			&m.KreditBermasalah,
			&m.TotalKredit,
			&m.CreatedAt,
			&m.CreatedBy,
			&m.UpdatedAt,
			&m.UpdatedBy,
		)
		if err != nil {
			return r, totalCount, pageCount, err
		}
		r = append(r, m)
	}

	pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))

	return r, totalCount, pageCount, err
}

func (c *KinerjaController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r others.KinerjaModel, err error) {
	r = others.KinerjaModel{}
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT
		trn_kinerja.id_kinerja,
		trn_kinerja.id_bumd,
		trn_kinerja.tahun_kinerja,
		trn_kinerja.ebit_kinerja,
		trn_kinerja.ebitda_kinerja,
		trn_kinerja.modal_sendiri_kinerja,
		trn_kinerja.penyusutan_kinerja,
		trn_kinerja.capital_employed_kinerja,
		trn_kinerja.total_aset_awal_kinerja,
		trn_kinerja.total_aset_akhir_kinerja,
		trn_kinerja.kas_kinerja,
		trn_kinerja.setara_kas_kinerja,
		trn_kinerja.kewajiban_lancar_kinerja,
		trn_kinerja.harta_lancar_kinerja,
		trn_kinerja.penjualan_bersih_kinerja,
		trn_kinerja.piutang_dagang_kinerja,
		trn_kinerja.harga_pokok_penjualan_kinerja,
		trn_kinerja.persediaan_kinerja,
		trn_kinerja.aktiva_tetap_kinerja,
		trn_kinerja.akumulasi_depresiasi_kinerja,
		trn_kinerja.kredit_bermasalah_kinerja,
		trn_kinerja.total_kredit_kinerja,
		trn_kinerja.created_at,
		trn_kinerja.created_by,
		trn_kinerja.updated_at,
		trn_kinerja.updated_by
	FROM trn_kinerja
	WHERE trn_kinerja.deleted_by = 0
	AND trn_kinerja.id_bumd = $1
	AND trn_kinerja.id_kinerja = $2
	`
	err = c.pgxConn.QueryRow(fCtx, q, idBumd, id).Scan(&r.IdKinerja, &r.IdBumd, &r.Tahun, &r.Ebit, &r.Ebitda, &r.ModalSendiri, &r.Penyusutan, &r.CapitalEmployed, &r.TotalAsetAwal, &r.TotalAsetAkhir, &r.Kas, &r.SetaraKas, &r.KewajibanLancar, &r.Hartalancar, &r.PenjualanBersih, &r.PiutangDagang, &r.HargaPokokPenjualan, &r.Persediaan, &r.AktivaTetap, &r.AkumulasiDepresiasi, &r.KreditBermasalah, &r.TotalKredit, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return r, err
	}

	return r, err
}

func (c *KinerjaController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *others.KinerjaForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO trn_kinerja (
		id_kinerja,
		id_bumd,
		tahun_kinerja,
		ebit_kinerja,
		ebitda_kinerja,
		modal_sendiri_kinerja,
		penyusutan_kinerja,
		capital_employed_kinerja,
		total_aset_awal_kinerja,
		total_aset_akhir_kinerja,
		kas_kinerja,
		setara_kas_kinerja,
		kewajiban_lancar_kinerja,
		harta_lancar_kinerja,
		penjualan_bersih_kinerja,
		piutang_dagang_kinerja,
		harga_pokok_penjualan_kinerja,
		persediaan_kinerja,
		aktiva_tetap_kinerja,
		akumulasi_depresiasi_kinerja,
		kredit_bermasalah_kinerja,
		total_kredit_kinerja,
		created_by
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23
	)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id KINERJA. - " + err.Error(),
		}
	}

	_, err = c.pgxConn.Exec(fCtx, q,
		id,
		idBumd,
		payload.Tahun,
		payload.Ebit,
		payload.Ebitda,
		payload.ModalSendiri,
		payload.Penyusutan,
		payload.CapitalEmployed,
		payload.TotalAsetAwal,
		payload.TotalAsetAkhir,
		payload.Kas,
		payload.SetaraKas,
		payload.KewajibanLancar,
		payload.Hartalancar,
		payload.PenjualanBersih,
		payload.PiutangDagang,
		payload.HargaPokokPenjualan,
		payload.Persediaan,
		payload.AktivaTetap,
		payload.AkumulasiDepresiasi,
		payload.KreditBermasalah,
		payload.TotalKredit,
		idUser,
	)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *KinerjaController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *others.KinerjaForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_kinerja SET
		tahun_kinerja = $1,
		ebit_kinerja = $2,
		ebitda_kinerja = $3,
		modal_sendiri_kinerja = $4,
		penyusutan_kinerja = $5,
		capital_employed_kinerja = $6,
		total_aset_awal_kinerja = $7,
		total_aset_akhir_kinerja = $8,
		kas_kinerja = $9,
		setara_kas_kinerja = $10,
		kewajiban_lancar_kinerja = $11,
		harta_lancar_kinerja = $12,
		penjualan_bersih_kinerja = $13,
		piutang_dagang_kinerja = $14,
		harga_pokok_penjualan_kinerja = $15,
		persediaan_kinerja = $16,
		aktiva_tetap_kinerja = $17,
		akumulasi_depresiasi_kinerja = $18,
		kredit_bermasalah_kinerja = $19,
		total_kredit_kinerja = $20,
		updated_by = $21
	WHERE id_bumd = $22 AND id_kinerja = $23
	`
	_, err = c.pgxConn.Exec(fCtx, q,
		payload.Tahun,
		payload.Ebit,
		payload.Ebitda,
		payload.ModalSendiri,
		payload.Penyusutan,
		payload.CapitalEmployed,
		payload.TotalAsetAwal,
		payload.TotalAsetAkhir,
		payload.Kas,
		payload.SetaraKas,
		payload.KewajibanLancar,
		payload.Hartalancar,
		payload.PenjualanBersih,
		payload.PiutangDagang,
		payload.HargaPokokPenjualan,
		payload.Persediaan,
		payload.AktivaTetap,
		payload.AkumulasiDepresiasi,
		payload.KreditBermasalah,
		payload.TotalKredit,
		idUser,
		idBumd,
		id,
	)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *KinerjaController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_kinerja SET
		deleted_by = $1,
		deleted_at = NOW()
	WHERE id_bumd = $2 AND id_kinerja = $3
	`
	_, err = c.pgxConn.Exec(fCtx, q, idUser, idBumd, id)
	if err != nil {
		return false, err
	}

	return true, err
}
