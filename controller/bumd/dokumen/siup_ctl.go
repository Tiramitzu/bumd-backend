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

type SiupController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewSiupController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *SiupController {
	return &SiupController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *SiupController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit,
	isSeumurHidup int,
	search string,
) (r []dokumen.SiupModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SIUP. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.SiupModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_siup WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_siup, nomor_siup, instansi_pemberi_siup, tanggal_siup, kualifikasi_siup, klasifikasi_siup, masa_berlaku_siup, file_siup, id_bumd,
	CASE
		WHEN masa_berlaku_siup IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM trn_siup
	WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nomor_siup ILIKE $%d OR instansi_pemberi_siup ILIKE $%d OR tanggal_siup ILIKE $%d OR klasifikasi_siup ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nomor_siup ILIKE $%d OR instansi_pemberi_siup ILIKE $%d OR tanggal_siup ILIKE $%d OR klasifikasi_siup ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}
	if isSeumurHidup != 0 {
		qCount += ` AND masa_berlaku_siup IS NULL`
		q += ` AND masa_berlaku_siup IS NULL`
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data SIUP. - " + err.Error(),
		}
	}

	q += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SIUP. - " + err.Error(),
		}
	}
	defer rows.Close()
	for rows.Next() {
		var m dokumen.SiupModel
		err = rows.Scan(&m.Id, &m.Nomor, &m.InstansiPemberi, &m.Tanggal, &m.Kualifikasi, &m.Klasifikasi, &m.MasaBerlaku, &m.File, &m.IdBumd, &m.IsSeumurHidup)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data SIUP. - " + err.Error(),
			}
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *SiupController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r dokumen.SiupModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SIUP. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_siup, nomor_siup, instansi_pemberi_siup, tanggal_siup, kualifikasi_siup, klasifikasi_siup, masa_berlaku_siup, file_siup, id_bumd,
	CASE
		WHEN masa_berlaku_siup IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM trn_siup
	WHERE id_siup = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.Id, &r.Nomor, &r.InstansiPemberi, &r.Tanggal, &r.Kualifikasi, &r.Klasifikasi, &r.MasaBerlaku, &r.File, &r.IdBumd, &r.IsSeumurHidup)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data SIUP tidak ditemukan",
			}
		}
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SIUP. - " + err.Error(),
		}
	}

	return r, err
}

func (c *SiupController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *dokumen.SiupForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))

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
	INSERT INTO trn_siup (id_siup, nomor_siup, instansi_pemberi_siup, tanggal_siup, kualifikasi_siup, klasifikasi_siup, id_bumd, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{Code: fasthttp.StatusInternalServerError, Message: "gagal membuat id SIUP. - " + err.Error()}
	}

	_, err = tx.Exec(context.Background(), q, id, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, idBumd, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data SIUP. - " + err.Error(),
		}
	}

	if payload.MasaBerlaku != nil {
		q = `UPDATE trn_siup SET masa_berlaku_siup = $1 WHERE id_siup = $2 AND id_bumd = $3`

		_, err = tx.Exec(context.Background(), q, payload.MasaBerlaku, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate data SIUP. - " + err.Error(),
			}
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
		objectName := "trn_siup/" + fileName
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
		q = `UPDATE trn_siup SET file_siup=$1 WHERE id_siup=$2 AND id_bumd=$3`
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

func (c *SiupController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *dokumen.SiupForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))

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

	var args []interface{}
	q := `
	UPDATE trn_siup
	SET nomor_siup = $1, instansi_pemberi_siup = $2, tanggal_siup = $3, kualifikasi_siup = $4, klasifikasi_siup = $5, updated_by = $6, updated_at = NOW()
	WHERE id_siup = $7 AND id_bumd = $8
	`
	args = append(args, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate data SIUP. - " + err.Error(),
		}
	}

	if payload.MasaBerlaku != nil {
		q = `UPDATE trn_siup SET masa_berlaku_siup = $1 WHERE id_siup = $2 AND id_bumd = $3`

		_, err = tx.Exec(context.Background(), q, payload.MasaBerlaku, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate data SIUP. - " + err.Error(),
			}
		}
	} else {
		q = `UPDATE trn_siup SET masa_berlaku_siup='' WHERE id_siup=$1 AND id_bumd=$2`
		_, err = tx.Exec(context.Background(), q, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate data SIUP. - " + err.Error(),
			}
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
		objectName := "trn_siup/" + fileName
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
		q = `UPDATE trn_siup SET file_siup=$1 WHERE id_siup=$2 AND id_bumd=$3`
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

func (c *SiupController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SIUP. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_siup
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_siup = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghapus data SIUP. - " + err.Error(),
		}
	}
	return true, err
}
