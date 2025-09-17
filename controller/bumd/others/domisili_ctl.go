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

type DomisiliController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewDomisiliController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *DomisiliController {
	return &DomisiliController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *DomisiliController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit,
	isSeumurHidup int,
	search string,
) (r []others.DomisiliModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]others.DomisiliModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_domisili WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_domisili, nomor_domisili, instansi_pemberi_domisili, tanggal_domisili, kualifikasi_domisili, klasifikasi_domisili, masa_berlaku_domisili, file_domisili, id_bumd,
	CASE
		WHEN masa_berlaku_domisili IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM trn_domisili
	WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nomor_domisili ILIKE $%d OR instansi_pemberi_domisili ILIKE $%d OR tanggal_domisili ILIKE $%d OR klasifikasi_domisili ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nomor_domisili ILIKE $%d OR instansi_pemberi_domisili ILIKE $%d OR tanggal_domisili ILIKE $%d OR klasifikasi_domisili ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}
	if isSeumurHidup != 0 {
		qCount += ` AND masa_berlaku_domisili IS NULL`
		q += ` AND masa_berlaku_domisili IS NULL`
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data DOMISILI: %w", err)
	}

	q += fmt.Sprintf(` ORDER BY id_domisili DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m others.DomisiliModel
		err = rows.Scan(&m.Id, &m.Nomor, &m.InstansiPemberi, &m.Tanggal, &m.Kualifikasi, &m.Klasifikasi, &m.MasaBerlaku, &m.File, &m.IdBumd, &m.IsSeumurHidup)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data DOMISILI: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *DomisiliController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r others.DomisiliModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_domisili, nomor_domisili, instansi_pemberi_domisili, tanggal_domisili, kualifikasi_domisili, klasifikasi_domisili, masa_berlaku_domisili, file_domisili, id_bumd,
	CASE
		WHEN masa_berlaku_domisili IS NULL THEN 1
		ELSE 0
	END as is_seumur_hidup
	FROM trn_domisili
	WHERE id_domisili = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.Id, &r.Nomor, &r.InstansiPemberi, &r.Tanggal, &r.Kualifikasi, &r.Klasifikasi, &r.MasaBerlaku, &r.File, &r.IdBumd, &r.IsSeumurHidup)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data DOMISILI tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
	}

	return r, err
}

func (c *DomisiliController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *others.DomisiliForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
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
	INSERT INTO trn_domisili (id_domisili, nomor_domisili, instansi_pemberi_domisili, tanggal_domisili, kualifikasi_domisili, klasifikasi_domisili, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id DOMISILI. - " + err.Error(),
		}
	}

	_, err = tx.Exec(context.Background(), q, id, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, idBumd, idUser)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data DOMISILI. - " + err.Error(),
		}
		return false, err
	}

	if payload.MasaBerlaku != nil {
		q = `
		UPDATE trn_domisili
		SET masa_berlaku_domisili = $1
		WHERE id_domisili = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(context.Background(), q, payload.MasaBerlaku, id, idBumd)
		if err != nil {
			return false, err
		}
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
		objectName := "trn_domisili/" + fileName
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
		q = `UPDATE trn_domisili SET file_domisili=$1 WHERE id_domisili=$2 AND id_bumd=$3`
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

func (c *DomisiliController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *others.DomisiliForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
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
	UPDATE trn_domisili
	SET nomor_domisili = $1, instansi_pemberi_domisili = $2, tanggal_domisili = $3, kualifikasi_domisili = $4, klasifikasi_domisili = $5, updated_by = $6, updated_at = NOW()
	WHERE id_domisili = $7 AND id_bumd = $8
	`
	args = append(args, payload.Nomor, payload.InstansiPemberi, payload.Tanggal, payload.Kualifikasi, payload.Klasifikasi, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.MasaBerlaku != nil {
		q = `
		UPDATE trn_domisili
		SET masa_berlaku_domisili = $1
		WHERE id_domisili = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(context.Background(), q, payload.MasaBerlaku, id, idBumd)
		if err != nil {
			return false, err
		}
	} else {
		q = `UPDATE trn_domisili SET masa_berlaku_domisili='' WHERE id_domisili=$1 AND id_bumd=$2`
		_, err = tx.Exec(context.Background(), q, id, idBumd)
		if err != nil {
			return false, err
		}
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
		objectName := "trn_domisili/" + fileName
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
		q = `UPDATE trn_domisili SET file_domisili=$1 WHERE id_domisili=$2`
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

func (c *DomisiliController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, fmt.Errorf("gagal mengambil data DOMISILI: %w", err)
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}
	q := `
	UPDATE trn_domisili
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_domisili = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(context.Background(), q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}
	return true, err
}
