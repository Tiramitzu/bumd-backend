package keuangan

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/keuangan"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type LaporanKeuanganController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewLaporanKeuanganController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *LaporanKeuanganController {
	return &LaporanKeuanganController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *LaporanKeuanganController) Index(
	ctx context.Context,
	user *jwt.Token,
	idBumd uuid.UUID,
	page,
	limit int,
) (r []keuangan.LaporanKeuanganModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, err
	}
	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	r = make([]keuangan.LaporanKeuanganModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_laporan_keuangan WHERE deleted_by = 0 AND id_bumd = $1`

	q := `
	SELECT
		id_laporan_keuangan,
		id_bumd,
		tlk.id_jenis_laporan,
		tlk.id_jenis_laporan_item,
		mjl.uraian_jenis_laporan,
		mjl_item.uraian_jenis_laporan,
		tahun_laporan_keuangan,
		jumlah_laporan_keuangan,
		file_laporan_keuangan,
		created_at,
		created_by,
		updated_at, 
		updated_by 
	FROM trn_laporan_keuangan tlk
	LEFT JOIN m_jenis_laporan mjl ON mjl.id_jenis_laporan = tlk.id_jenis_laporan
	LEFT JOIN m_jenis_laporan mjl_item ON mjl_item.id_jenis_laporan = tlk.id_jenis_laporan_item
	WHERE id_bumd = $1 
		AND deleted_by = 0 
	LIMIT $2 
	OFFSET $3
	`

	err = c.pgxConn.QueryRow(ctx, qCount, idBumd).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data Laporan Keuangan. - " + err.Error(),
		}
	}

	rows, err := c.pgxConn.Query(ctx, q, idBumd, limit, offset)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		var m keuangan.LaporanKeuanganModel
		err = rows.Scan(&m.Id, &m.IdBumd, &m.IdJenisLaporan, &m.IdJenisLaporanItem, &m.NamaJenisLaporan, &m.NamaJenisLaporanItem, &m.Tahun, &m.Jumlah, &m.File, &m.CreatedAt, &m.CreatedBy, &m.UpdatedAt, &m.UpdatedBy)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
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

func (c *LaporanKeuanganController) View(ctx context.Context, user *jwt.Token, idBumd uuid.UUID, id uuid.UUID) (r keuangan.LaporanKeuanganModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
	}
	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT 
		id_laporan_keuangan,
		id_bumd,
		tlk.id_jenis_laporan,
		tlk.id_jenis_laporan_item,
		mjl.uraian_jenis_laporan,
		mjl_item.uraian_jenis_laporan,
		tahun_laporan_keuangan,
		jumlah_laporan_keuangan,
		file_laporan_keuangan,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM trn_laporan_keuangan tlk
	LEFT JOIN m_jenis_laporan mjl ON mjl.id_jenis_laporan = tlk.id_jenis_laporan
	LEFT JOIN m_jenis_laporan mjl_item ON mjl_item.id_jenis_laporan = tlk.id_jenis_laporan_item
	WHERE id_bumd = $1 AND id_laporan_keuangan = $2 AND deleted_by = 0`

	err = c.pgxConn.QueryRow(ctx, q, idBumd, id).Scan(&r.Id, &r.IdBumd, &r.IdJenisLaporan, &r.IdJenisLaporanItem, &r.NamaJenisLaporan, &r.NamaJenisLaporanItem, &r.Tahun, &r.Jumlah, &r.File, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
	}
	return r, err
}

func (c *LaporanKeuanganController) Create(ctx context.Context, user *jwt.Token, idBumd uuid.UUID, payload *keuangan.LaporanKeuanganForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
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
	INSERT INTO trn_laporan_keuangan (
		id_laporan_keuangan,
		id_bumd,
		id_jenis_laporan,
		id_jenis_laporan_item,
		tahun_laporan_keuangan,
		jumlah_laporan_keuangan,
		created_by,
		created_at
	) VALUES (
		$1,
		$2,
		$3,
		$4, 
		$5,
		$6,
		$7,
		NOW()
	)`

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id Laporan Keuangan. - " + err.Error(),
		}
	}

	_, err = tx.Exec(ctx, q, id, idBumd, payload.IdJenisLaporan, payload.IdJenisLaporanItem, payload.Tahun, payload.Jumlah, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat data Laporan Keuangan. - " + err.Error(),
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
		objectName := "trn_laporan_keuangan/" + fileName
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
		q = `UPDATE trn_laporan_keuangan SET file_laporan_keuangan=$1 WHERE id_laporan_keuangan=$2 AND id_bumd=$3`
		_, err = tx.Exec(ctx, q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *LaporanKeuanganController) Update(ctx context.Context, user *jwt.Token, idBumd uuid.UUID, id uuid.UUID, payload *keuangan.LaporanKeuanganForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
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

	q := `UPDATE trn_laporan_keuangan SET id_jenis_laporan=$1, id_jenis_laporan_item=$2, tahun_laporan_keuangan=$3, jumlah_laporan_keuangan=$4, updated_by=$5, updated_at=NOW() WHERE id_laporan_keuangan=$6 AND id_bumd=$7`
	_, err = tx.Exec(ctx, q, payload.IdJenisLaporan, payload.IdJenisLaporanItem, payload.Tahun, payload.Jumlah, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate data Laporan Keuangan. - " + err.Error(),
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
		objectName := "trn_laporan_keuangan/" + fileName
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
		q = `UPDATE trn_laporan_keuangan SET file_laporan_keuangan=$1 WHERE id_laporan_keuangan=$2 AND id_bumd=$3`
		_, err = tx.Exec(ctx, q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *LaporanKeuanganController) Delete(ctx context.Context, user *jwt.Token, idBumd uuid.UUID, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Laporan Keuangan. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `UPDATE trn_laporan_keuangan SET deleted_by=$1, deleted_at=NOW() WHERE id_laporan_keuangan=$2 AND id_bumd=$3`
	_, err = c.pgxConn.Exec(ctx, q, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghapus data Laporan Keuangan. - " + err.Error(),
		}
	}

	return true, err
}
