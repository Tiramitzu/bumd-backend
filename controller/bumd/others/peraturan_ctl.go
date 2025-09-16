package others

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type PeraturanController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewPeraturanController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *PeraturanController {
	return &PeraturanController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *PeraturanController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit int,
	search string,
) (r []others.PeraturanModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]others.PeraturanModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_peraturan WHERE trn_peraturan.deleted_by = 0 AND trn_peraturan.id_bumd = $1`
	q := `
	SELECT id_peraturan, nomor_peraturan, tanggal_berlaku_peraturan, keterangan_peraturan, file_peraturan, id_bumd, jenis_peraturan, m_jenis_dokumen.nama_jd
	FROM trn_peraturan
	LEFT JOIN m_jenis_dokumen ON trn_peraturan.jenis_peraturan = m_jenis_dokumen.id_jd
	WHERE trn_peraturan.deleted_by = 0 AND trn_peraturan.id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nomor_peraturan ILIKE $%d OR tanggal_berlaku_peraturan ILIKE $%d OR keterangan_peraturan ILIKE $%d OR nama_jenis_peraturan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nomor_peraturan ILIKE $%d OR tanggal_berlaku_peraturan ILIKE $%d OR keterangan_peraturan ILIKE $%d OR nama_jenis_peraturan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data PERATURAN: %w", err)
	}

	q += fmt.Sprintf(` ORDER BY trn_peraturan.created_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data PERATURAN: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m others.PeraturanModel
		err = rows.Scan(&m.Id, &m.Nomor, &m.TanggalBerlaku, &m.KeteranganPeraturan, &m.FilePeraturan, &m.IdBumd, &m.JenisPeraturan, &m.NamaJenisPeraturan)
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

func (c *PeraturanController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r others.PeraturanModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_peraturan, nomor_peraturan, tanggal_berlaku_peraturan, keterangan_peraturan, file_peraturan, id_bumd, jenis_peraturan, m_jenis_dokumen.nama_jd
	FROM trn_peraturan
	LEFT JOIN m_jenis_dokumen ON trn_peraturan.jenis_peraturan = m_jenis_dokumen.id_jd
	WHERE trn_peraturan.id_peraturan = $1 AND trn_peraturan.id_bumd = $2 AND trn_peraturan.deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.Id, &r.Nomor, &r.TanggalBerlaku, &r.KeteranganPeraturan, &r.FilePeraturan, &r.IdBumd, &r.JenisPeraturan, &r.NamaJenisPeraturan)
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

func (c *PeraturanController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *others.PeraturanForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

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

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO trn_peraturan (id_peraturan, nomor_peraturan, keterangan_peraturan, id_bumd, jenis_peraturan, created_by) VALUES ($1, $2, $3, $4, $5, $6)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id PERATURAN. - " + err.Error(),
		}
	}

	_, err = tx.Exec(context.Background(), q, id, payload.Nomor, payload.KeteranganPeraturan, idBumd, payload.JenisPeraturan, idUser)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data PERATURAN. - " + err.Error(),
		}
		return false, err
	}

	if payload.TanggalBerlaku != nil {
		q = `
		UPDATE trn_peraturan
		SET tanggal_berlaku_peraturan = $1
		WHERE id_peraturan = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(context.Background(), q, payload.TanggalBerlaku, id, idBumd)
		if err != nil {
			return false, err
		}
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
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.FilePeraturan.Size,
			minio.PutObjectOptions{ContentType: payload.FilePeraturan.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q = `UPDATE trn_peraturan SET file_peraturan=$1 WHERE id_peraturan=$2 AND id_bumd=$3`
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

func (c *PeraturanController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *others.PeraturanForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

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

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	var args []interface{}
	q := `
	UPDATE trn_peraturan
	SET nomor_peraturan = $1, keterangan_peraturan = $2, jenis_peraturan = $3, updated_by = $4, updated_at = NOW()
	WHERE id_peraturan = $5 AND id_bumd = $6
	`
	args = append(args, payload.Nomor, payload.KeteranganPeraturan, payload.JenisPeraturan, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.TanggalBerlaku != nil {
		q = `
		UPDATE trn_peraturan
		SET tanggal_berlaku_peraturan = $1
		WHERE id_peraturan = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(context.Background(), q, payload.TanggalBerlaku, id, idBumd)
		if err != nil {
			return false, err
		}
	} else {
		q = `UPDATE trn_peraturan SET tanggal_berlaku_peraturan=NULL WHERE id_peraturan=$2 AND id_bumd=$3`
		_, err = tx.Exec(context.Background(), q, id, idBumd)
		if err != nil {
			return false, err
		}
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
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.FilePeraturan.Size,
			minio.PutObjectOptions{ContentType: payload.FilePeraturan.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q = `UPDATE trn_peraturan SET file_peraturan=$1 WHERE id_peraturan=$2`
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

func (c *PeraturanController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
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
	UPDATE trn_peraturan
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_peraturan = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(context.Background(), q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}
	return true, err
}
