package dokumen

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/dokumen"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type PerdaPendirianController struct {
	pgxConn *pgxpool.Pool
}

func NewPerdaPendirianController(pgxConn *pgxpool.Pool) *PerdaPendirianController {
	return &PerdaPendirianController{pgxConn: pgxConn}
}

func (c *PerdaPendirianController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit int,
	search string,
	modalDasarMin,
	modalDasarMax float64,
) (r []dokumen.PerdaPendirianModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.PerdaPendirianModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM dkmn_perda_pendirian WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, nomor, tanggal, keterangan, file, modal_dasar FROM dkmn_perda_pendirian WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(
			` AND (nomor ILIKE $%d OR tanggal ILIKE $%d OR keterangan ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		q += fmt.Sprintf(
			` AND (nomor ILIKE $%d OR tanggal ILIKE $%d OR keterangan ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}
	if modalDasarMin > 0 {
		qCount += fmt.Sprintf(` AND modal_dasar >= $%d`, len(args)+1)
		q += fmt.Sprintf(` AND modal_dasar >= $%d`, len(args)+1)
		args = append(args, modalDasarMin)
	}
	if modalDasarMax > 0 {
		qCount += fmt.Sprintf(` AND modal_dasar <= $%d`, len(args)+1)
		q += fmt.Sprintf(` AND modal_dasar <= $%d`, len(args)+1)
		args = append(args, modalDasarMax)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Akta Notaris: %w", err)
	}

	q += fmt.Sprintf(`ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m dokumen.PerdaPendirianModel
		err = rows.Scan(&m.ID, &m.Nomor, &m.Tanggal, &m.Keterangan, &m.File, &m.ModalDasar)
		m.IDBumd = int32(idBumd)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Bentuk Usaha: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *PerdaPendirianController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r dokumen.PerdaPendirianModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, nomor, tanggal, keterangan, file, modal_dasar FROM dkmn_perda_pendirian WHERE id = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.ID, &r.Nomor, &r.Tanggal, &r.Keterangan, &r.File, &r.ModalDasar)
	r.IDBumd = int32(idBumd)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data Perda Pendirian tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
	}

	return r, err
}

func (c *PerdaPendirianController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *dokumen.PerdaPendirianForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	q := `
	INSERT INTO dkmn_perda_pendirian (nomor, tanggal, keterangan, modal_dasar, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	var id int
	err = tx.QueryRow(context.Background(), q, payload.Nomor, payload.Tanggal, payload.Keterangan, payload.ModalDasar, idBumd, idUser).Scan(&id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data Perda Pendirian. - " + err.Error(),
		}
	}

	if payload.File != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
		}
		defer src.Close()

		// upload file
		objectName := "dkmn_perda_pendirian/" + fileName

		// update file
		q = `UPDATE dkmn_perda_pendirian SET file=$1 WHERE id=$2 AND id_bumd=$3`
		_, err = tx.Exec(context.Background(), q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *PerdaPendirianController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int, payload *dokumen.PerdaPendirianForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	q := `
	UPDATE dkmn_perda_pendirian
	SET nomor = $1, tanggal = $2, keterangan = $3, modal_dasar = $4, updated_by = $5
	WHERE id = $6 AND id_bumd = $7
	`
	args = append(args, payload.Nomor, payload.Tanggal, payload.Keterangan, payload.ModalDasar, idUser, id, idBumd)

	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.File != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
		}
		defer src.Close()

		// upload file
		objectName := "dkmn_perda_pendirian/" + fileName

		// update file
		q = `UPDATE dkmn_perda_pendirian SET file=$1 WHERE id=$2 AND id_bumd=$3`
		_, err = tx.Exec(context.Background(), q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *PerdaPendirianController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE dkmn_perda_pendirian
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3 AND id_bumd = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id, idBumd)
	if err != nil {
		return false, err
	}

	return true, err
}
