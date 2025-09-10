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

type NibController struct {
	pgxConn *pgxpool.Pool
}

func NewNibController(pgxConn *pgxpool.Pool) *NibController {
	return &NibController{pgxConn: pgxConn}
}

func (c *NibController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit,
	isSeumurHidup int,
	search string,
) (r []dokumen.NibModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.NibModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM dkmn_nib WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, nomor, instansi_pemberi, tanggal, kualifikasi, klasifikasi, masa_berlaku, file, id_bumd,
	CASE
		WHEN masa_berlaku IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM dkmn_nib
	WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nomor ILIKE $%d OR instansi_pemberi ILIKE $%d OR tanggal ILIKE $%d OR klasifikasi ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nomor ILIKE $%d OR instansi_pemberi ILIKE $%d OR tanggal ILIKE $%d OR klasifikasi ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}
	if isSeumurHidup != 0 {
		qCount += ` AND masa_berlaku IS NULL`
		q += ` AND masa_berlaku IS NULL`
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data NIB: %w", err)
	}

	q += fmt.Sprintf(` ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data NIB: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m dokumen.NibModel
		err = rows.Scan(&m.ID, &m.Nomor, &m.InstansiPemberi, &m.Tanggal, &m.Kualifikasi, &m.Klasifikasi, &m.MasaBerlaku, &m.File, &m.IDBumd, &m.IsSeumurHidup)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data NIB: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *NibController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r dokumen.NibModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, nomor, instansi_pemberi, tanggal, kualifikasi, klasifikasi, masa_berlaku, file, id_bumd,
	CASE
		WHEN masa_berlaku IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM dkmn_nib
	WHERE id = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.ID, &r.Nomor, &r.InstansiPemberi, &r.Tanggal, &r.Kualifikasi, &r.Klasifikasi, &r.MasaBerlaku, &r.File, &r.IDBumd, &r.IsSeumurHidup)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data NIB tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data NIB: %w", err)
	}

	return r, err
}

func (c *NibController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *dokumen.NibForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
		return false, err
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

	q := `
	INSERT INTO dkmn_nib (nomor, instansi_pemberi, tanggal, kualifikasi, klasifikasi, masa_berlaku, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`

	var id int
	err = tx.QueryRow(context.Background(), q, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, payload.MasaBerlaku, idBumd, idUser).Scan(&id)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data NIB. - " + err.Error(),
		}
		return false, err
	}

	if payload.File != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "dkmn_nib/" + fileName

		// update file
		q = `UPDATE dkmn_nib SET file=$1 WHERE id=$2 AND id_bumd=$3`
		_, err = tx.Exec(context.Background(), q, objectName, id, idBumd)
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
			return false, err
		}
	}

	return true, err
}

func (c *NibController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int, payload *dokumen.NibForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
		return false, err
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
	UPDATE dkmn_nib
	SET nomor = $1, instansi_pemberi = $2, tanggal = $3, kualifikasi = $4, klasifikasi = $5, masa_berlaku = $6, updated_by = $7
	WHERE id = $8 AND id_bumd = $9
	`
	args = append(args, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, payload.MasaBerlaku, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.File != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "dkmn_nib/" + fileName

		// update file
		q = `UPDATE dkmn_nib SET file=$1 WHERE id=$2`
		_, err = tx.Exec(context.Background(), q, objectName, id)
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
			return false, err
		}
	}

	return true, err
}

func (c *NibController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}
	q := `
	UPDATE dkmn_nib
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3 AND id_bumd = $4
	`
	_, err = c.pgxConn.Exec(context.Background(), q, idUser, time.Now(), id, idBumd)
	if err != nil {
		return false, err
	}
	return true, err
}
