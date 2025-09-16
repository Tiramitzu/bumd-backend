package dokumen

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/dokumen"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type AktaNotarisController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewAktaNotarisController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *AktaNotarisController {
	return &AktaNotarisController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *AktaNotarisController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit int,
	search string,
) (r []dokumen.AktaNotarisModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.AktaNotarisModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_akta_notaris WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_akta_notaris, nomor_akta_notaris, notaris_akta_notaris, tanggal_akta_notaris, keterangan_akta_notaris, file_akta_notaris FROM trn_akta_notaris WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(
			` AND (nomor_akta_notaris ILIKE $%d OR tanggal_akta_notaris ILIKE $%d OR tanggal_akta_notaris ILIKE $%d OR keterangan_akta_notaris ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		q += fmt.Sprintf(
			` AND (nomor_akta_notaris ILIKE $%d OR tanggal_akta_notaris ILIKE $%d OR tanggal_akta_notaris ILIKE $%d OR keterangan_akta_notaris ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Akta Notaris: %w", err)
	}

	q += fmt.Sprintf(`ORDER BY id_akta_notaris DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m dokumen.AktaNotarisModel
		err = rows.Scan(&m.Id, &m.Nomor, &m.Notaris, &m.Tanggal, &m.Keterangan, &m.File)
		m.IdBumd = idBumd
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Akta Notaris: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *AktaNotarisController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r dokumen.AktaNotarisModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_akta_notaris, nomor_akta_notaris, notaris_akta_notaris, tanggal_akta_notaris, keterangan_akta_notaris, file_akta_notaris FROM trn_akta_notaris WHERE id_akta_notaris = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.Id, &r.Nomor, &r.Notaris, &r.Tanggal, &r.Keterangan, &r.File)
	r.IdBumd = idBumd
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data Akta Notaris tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	return r, err
}

func (c *AktaNotarisController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *dokumen.AktaNotarisForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
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

	q := `
	INSERT INTO trn_akta_notaris (id_akta_notaris, nomor_akta_notaris, notaris_akta_notaris, tanggal_akta_notaris, keterangan_akta_notaris, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_akta_notaris
	`

	var id uuid.UUID
	id, err = uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id Akta Notaris. - " + err.Error(),
		}
	}
	_, err = tx.Exec(context.Background(), q, id, payload.Nomor, payload.Notaris, payload.Tanggal, payload.Keterangan, idBumd, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data Akta Notaris. - " + err.Error(),
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
		objectName := "trn_akta_notaris/" + fileName
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
		q = `UPDATE trn_akta_notaris SET file_akta_notaris=$1 WHERE id_akta_notaris=$2 AND id_bumd=$3`
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

func (c *AktaNotarisController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *dokumen.AktaNotarisForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
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

	var args []interface{}
	q := `
	UPDATE trn_akta_notaris
	SET nomor_akta_notaris = $1, notaris_akta_notaris = $2, tanggal_akta_notaris = $3, keterangan_akta_notaris = $4, updated_by = $5, updated_at = NOW()
	WHERE id_akta_notaris = $6 AND id_bumd = $7
	`
	args = append(args, payload.Nomor, payload.Notaris, payload.Tanggal, payload.Keterangan, idUser, id, idBumd)

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
		objectName := "trn_akta_notaris/" + fileName

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
		q = `UPDATE trn_akta_notaris SET file_akta_notaris=$1 WHERE id_akta_notaris=$2 AND id_bumd=$3`
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

func (c *AktaNotarisController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_akta_notaris
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_akta_notaris = $2 AND id_bumd = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}

	return true, err
}
