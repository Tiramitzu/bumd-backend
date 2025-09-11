package others

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type PeraturanController struct {
	pgxConn *pgxpool.Pool
}

func NewPeraturanController(pgxConn *pgxpool.Pool) *PeraturanController {
	return &PeraturanController{pgxConn: pgxConn}
}

func (c *PeraturanController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit int,
	search string,
) (r []others.PeraturanModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]others.PeraturanModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_peraturan WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, nomor, tanggal_berlaku, keterangan_peraturan, file_peraturan, id_bumd, jenis_peraturan, nama_jenis_peraturan
	FROM trn_peraturan
	WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nomor ILIKE $%d OR tanggal_berlaku ILIKE $%d OR keterangan_peraturan ILIKE $%d OR nama_jenis_peraturan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nomor ILIKE $%d OR tanggal_berlaku ILIKE $%d OR keterangan_peraturan ILIKE $%d OR nama_jenis_peraturan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data PERATURAN: %w", err)
	}

	q += fmt.Sprintf(` ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data PERATURAN: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m others.PeraturanModel
		err = rows.Scan(&m.ID, &m.Nomor, &m.TanggalBerlaku, &m.KeteranganPeraturan, &m.FilePeraturan, &m.IDBumd, &m.JenisPeraturan, &m.NamaJenisPeraturan)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data PERATURAN: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *PeraturanController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string) (r others.PeraturanModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, nomor, tanggal_berlaku, keterangan_peraturan, file_peraturan, id_bumd, jenis_peraturan, nama_jenis_peraturan
	FROM trn_peraturan
	WHERE id = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.ID, &r.Nomor, &r.TanggalBerlaku, &r.KeteranganPeraturan, &r.FilePeraturan, &r.IDBumd, &r.JenisPeraturan, &r.NamaJenisPeraturan)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data PERATURAN tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data PERATURAN: %w", err)
	}

	return r, err
}

func (c *PeraturanController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *others.PeraturanForm) (r bool, err error) {
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
	INSERT INTO trn_peraturan (nomor, tanggal_berlaku, keterangan_peraturan, file_peraturan, id_bumd, jenis_peraturan, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	var id int
	err = tx.QueryRow(context.Background(), q, payload.Nomor, payload.TanggalBerlaku, payload.KeteranganPeraturan, payload.FilePeraturan, idBumd, payload.JenisPeraturan, idUser).Scan(&id)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data PERATURAN. - " + err.Error(),
		}
		return false, err
	}

	if payload.FilePeraturan != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.FilePeraturan.Filename)

		src, err := payload.FilePeraturan.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "trn_peraturan/" + fileName

		// update file
		q = `UPDATE trn_peraturan SET file_peraturan=$1 WHERE id=$2 AND id_bumd=$3`
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

func (c *PeraturanController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string, payload *others.PeraturanForm) (r bool, err error) {
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
	UPDATE trn_peraturan
	SET nomor = $1, tanggal_berlaku = $2, keterangan_peraturan = $3, file_peraturan = $4, jenis_peraturan = $5, updated_by = $6, updated_at = NOW()
	WHERE id = $8 AND id_bumd = $9
	`
	args = append(args, payload.Nomor, payload.TanggalBerlaku, payload.KeteranganPeraturan, payload.FilePeraturan, payload.JenisPeraturan, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.FilePeraturan != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.FilePeraturan.Filename)

		src, err := payload.FilePeraturan.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "trn_peraturan/" + fileName

		// update file
		q = `UPDATE trn_peraturan SET file_peraturan=$1 WHERE id=$2`
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

func (c *PeraturanController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}
	q := `
	UPDATE trn_peraturan
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(context.Background(), q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}
	return true, err
}
