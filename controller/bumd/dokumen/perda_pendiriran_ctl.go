package dokumen

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/dokumen"
	"microdata/kemendagri/bumd/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type PerdaPendirianController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewPerdaPendirianController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *PerdaPendirianController {
	return &PerdaPendirianController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *PerdaPendirianController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit int,
	search string,
	modalDasarMin,
	modalDasarMax float64,
) (r []dokumen.PerdaPendirianModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
	}
	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.PerdaPendirianModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_perda_pendirian WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_perda_pendirian, nomor_perda_pendirian, tanggal_perda_pendirian, keterangan_perda_pendirian, file_perda_pendirian, modal_dasar_perda_pendirian FROM trn_perda_pendirian WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(
			` AND (nomor_perda_pendirian ILIKE $%d OR tanggal_perda_pendirian ILIKE $%d OR keterangan_perda_pendirian ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		q += fmt.Sprintf(
			` AND (nomor_perda_pendirian ILIKE $%d OR tanggal_perda_pendirian ILIKE $%d OR keterangan_perda_pendirian ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}
	if modalDasarMin > 0 {
		qCount += fmt.Sprintf(` AND modal_dasar_perda_pendirian >= $%d`, len(args)+1)
		q += fmt.Sprintf(` AND modal_dasar_perda_pendirian >= $%d`, len(args)+1)
		args = append(args, modalDasarMin)
	}
	if modalDasarMax > 0 {
		qCount += fmt.Sprintf(` AND modal_dasar_perda_pendirian <= $%d`, len(args)+1)
		q += fmt.Sprintf(` AND modal_dasar_perda_pendirian <= $%d`, len(args)+1)
		args = append(args, modalDasarMax)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Akta Notaris: %w", err)
	}

	q += fmt.Sprintf(`ORDER BY id_perda_pendirian DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m dokumen.PerdaPendirianModel
		err = rows.Scan(&m.Id, &m.Nomor, &m.Tanggal, &m.Keterangan, &m.File, &m.ModalDasar)
		m.IdBumd = idBumd
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

func (c *PerdaPendirianController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r dokumen.PerdaPendirianModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_perda_pendirian, nomor_perda_pendirian, tanggal_perda_pendirian, keterangan_perda_pendirian, file_perda_pendirian, modal_dasar_perda_pendirian FROM trn_perda_pendirian WHERE id_perda_pendirian = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	r.IdBumd = idBumd
	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.Id, &r.Nomor, &r.Tanggal, &r.Keterangan, &r.File, &r.ModalDasar)
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

func (c *PerdaPendirianController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *dokumen.PerdaPendirianForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
	}

	if idBumdClaims != uuid.Nil {
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

	// Parse modal_dasar from string to float64
	modalDasar, err := strconv.ParseFloat(payload.ModalDasar, 64)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Modal dasar harus berupa angka yang valid",
		}
	}

	q := `
	INSERT INTO trn_perda_pendirian (id_perda_pendirian, nomor_perda_pendirian, tanggal_perda_pendirian, keterangan_perda_pendirian, modal_dasar_perda_pendirian, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_perda_pendirian
	`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id Perda Pendirian. - " + err.Error(),
		}
	}

	_, err = tx.Exec(context.Background(), q, id, payload.Nomor, payload.Tanggal, payload.Keterangan, modalDasar, idBumd, idUser)
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
		objectName := "trn_perda_pendirian/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.File.Size,
			minio.PutObjectOptions{ContentType: payload.File.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q = `UPDATE trn_perda_pendirian SET file_perda_pendirian=$1 WHERE id_perda_pendirian=$2 AND id_bumd=$3`
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

func (c *PerdaPendirianController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *dokumen.PerdaPendirianForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
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

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	// Parse modal_dasar from string to float64
	modalDasar, err := strconv.ParseFloat(payload.ModalDasar, 64)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Modal dasar harus berupa angka yang valid",
		}
	}

	var args []interface{}
	q := `
	UPDATE trn_perda_pendirian
	SET nomor_perda_pendirian = $1, tanggal_perda_pendirian = $2, keterangan_perda_pendirian = $3, modal_dasar_perda_pendirian = $4, updated_by = $5, updated_at = NOW()
	WHERE id_perda_pendirian = $6 AND id_bumd = $7
	`
	args = append(args, payload.Nomor, payload.Tanggal, payload.Keterangan, modalDasar, idUser, id, idBumd)

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
		objectName := "trn_perda_pendirian/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.File.Size,
			minio.PutObjectOptions{ContentType: payload.File.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q = `UPDATE trn_perda_pendirian SET file_perda_pendirian=$1 WHERE id_perda_pendirian=$2 AND id_bumd=$3`
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

func (c *PerdaPendirianController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Perda Pendirian: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_perda_pendirian
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_perda_pendirian = $2 AND id_bumd = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}

	return true, err
}
