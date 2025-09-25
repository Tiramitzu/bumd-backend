package bumd

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd"
	"microdata/kemendagri/bumd/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type BumdController struct {
	pgxConn        *pgxpool.Pool
	pgxConnMstData *pgxpool.Pool
	minioConn      *utils.MinioConn
}

func NewBumdController(pgxConn *pgxpool.Pool, pgxConnMstData *pgxpool.Pool, minioConn *utils.MinioConn) *BumdController {
	return &BumdController{pgxConn: pgxConn, pgxConnMstData: pgxConnMstData, minioConn: minioConn}
}

func (c *BumdController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	page, limit int,
	nama string,
	penerapanSPI bool,
	indukPerusahaan uuid.UUID,
) (r []bumd.BumdModel, totalCount, pageCount int, err error) {
	r = make([]bumd.BumdModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))
	offset := limit * (page - 1)

	var args []interface{}
	qCount := `
	SELECT COALESCE(COUNT(*), 0) FROM trn_bumd WHERE deleted_by = 0
	`

	q := `
	WITH t_induk_perusahaan AS (
		SELECT id_bumd, nama_bumd as nama_induk_perusahaan
		FROM trn_bumd
		WHERE deleted_by = 0
	)
	SELECT
		b.id_bumd,
		b.id_daerah,
		b.id_bentuk_hukum,
		b.id_bentuk_usaha,
		b.id_induk_perusahaan,
		COALESCE(t_induk_perusahaan.nama_induk_perusahaan, '-') as nama_induk_perusahaan,
		b.penerapan_spi_bumd,
		mbbh.nama_bbh as bentuk_badan_hukum,
		mbbu.nama_bu as bentuk_usaha,
		b.nama_bumd,
		b.deskripsi_bumd,
		b.alamat_bumd,
		b.no_telp_bumd,
		b.no_fax_bumd,
		b.email_bumd,
		b.website_bumd,
		b.narahubung_bumd,
		b.npwp_bumd,
		b.npwp_pemberi_bumd,
		b.npwp_file_bumd,
		b.file_spi_bumd,
		b.logo_bumd,
		b.created_at,
		b.created_by,
		b.updated_at,
		b.updated_by
	FROM trn_bumd b
		LEFT JOIN m_bentuk_badan_hukum mbbh ON mbbh.id_bbh = b.id_bentuk_hukum
		LEFT JOIN m_bentuk_usaha mbbu ON mbbu.id_bu = b.id_bentuk_usaha
		LEFT JOIN t_induk_perusahaan ON t_induk_perusahaan.id_bumd = b.id_induk_perusahaan
	WHERE b.deleted_by = 0
	`
	if idDaerah > 0 {
		qCount += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}
	if nama != "" {
		qCount += fmt.Sprintf(` AND nama_bumd ILIKE $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.nama_bumd ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}
	if penerapanSPI {
		qCount += fmt.Sprintf(` AND penerapan_spi_bumd = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.penerapan_spi_bumd = $%d`, len(args)+1)
		args = append(args, penerapanSPI)
	}
	if indukPerusahaan != uuid.Nil {
		qCount += fmt.Sprintf(` AND id_induk_perusahaan = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.id_induk_perusahaan = $%d`, len(args)+1)
		args = append(args, indukPerusahaan)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data BUMD. - " + err.Error(),
		}
	}

	q += fmt.Sprintf(`
	LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data BUMD. - " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m bumd.BumdModel
		err = rows.Scan(
			&m.Id,
			&m.IdDaerah,
			&m.IdBentukHukum,
			&m.IdBentukUsaha,
			&m.IdIndukPerusahaan,
			&m.NamaIndukPerusahaan,
			&m.PenerapanSPI,
			&m.BentukBadanHukum,
			&m.BentukUsaha,
			&m.Nama,
			&m.Deskripsi,
			&m.Alamat,
			&m.NoTelp,
			&m.NoFax,
			&m.Email,
			&m.Website,
			&m.Narahubung,
			&m.NPWP,
			&m.NPWPPemberi,
			&m.NPWPFile,
			&m.SPIFile,
			&m.Logo,
			&m.CreatedAt,
			&m.CreatedBy,
			&m.UpdatedAt,
			&m.UpdatedBy,
		)

		q = `SELECT nama_daerah, id_prop FROM data.m_daerah WHERE id_daerah = $1`
		err = c.pgxConnMstData.QueryRow(fCtx, q, m.IdDaerah).Scan(&m.NamaDaerah, &m.IdProvinsi)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data Daerah. - " + err.Error(),
			}
		}

		if m.IdDaerah != m.IdProvinsi {
			q = `SELECT nama_daerah FROM data.m_daerah WHERE id_daerah = $1`
			err = c.pgxConnMstData.QueryRow(fCtx, q, m.IdProvinsi).Scan(&m.NamaProvinsi)
			if err != nil {
				return r, totalCount, pageCount, utils.RequestError{
					Code:    fasthttp.StatusInternalServerError,
					Message: "gagal mengambil data Daerah. - " + err.Error(),
				}
			}
		}

		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data BUMD. - " + err.Error(),
			}
		}
		r = append(r, m)
	}

	// page info
	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *BumdController) View(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	id uuid.UUID,
) (r bumd.BumdModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	WITH t_induk_perusahaan AS (
		SELECT id_bumd, nama_bumd as nama_induk_perusahaan
		FROM trn_bumd
		WHERE deleted_by = 0
	)
	SELECT
		b.id_bumd,
		b.id_daerah,
		b.id_bentuk_hukum,
		b.id_bentuk_usaha,
		b.id_induk_perusahaan,
		COALESCE(t_induk_perusahaan.nama_induk_perusahaan, '-') as nama_induk_perusahaan,
		b.penerapan_spi_bumd,
		mbbh.nama_bbh as bentuk_badan_hukum,
		mbbu.nama_bu as bentuk_usaha,
		b.nama_bumd,
		b.deskripsi_bumd,
		b.alamat_bumd,
		b.no_telp_bumd,
		b.no_fax_bumd,
		b.email_bumd,
		b.website_bumd,
		b.narahubung_bumd,
		b.npwp_bumd,
		b.npwp_pemberi_bumd,
		b.npwp_file_bumd,
		b.file_spi_bumd,
		b.logo_bumd,
		b.created_at,
		b.created_by,
		b.updated_at,
		b.updated_by
	FROM trn_bumd b
		LEFT JOIN m_bentuk_badan_hukum mbbh ON mbbh.id_bbh = b.id_bentuk_hukum
		LEFT JOIN m_bentuk_usaha mbbu ON mbbu.id_bu = b.id_bentuk_usaha
		LEFT JOIN t_induk_perusahaan ON t_induk_perusahaan.id_bumd = b.id_induk_perusahaan
	WHERE b.id_bumd = $1 AND b.deleted_by = 0
	`
	args = append(args, id)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND b.id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(
		&r.Id,
		&r.IdDaerah,
		&r.IdBentukHukum,
		&r.IdBentukUsaha,
		&r.IdIndukPerusahaan,
		&r.NamaIndukPerusahaan,
		&r.PenerapanSPI,
		&r.BentukBadanHukum,
		&r.BentukUsaha,
		&r.Nama,
		&r.Deskripsi,
		&r.Alamat,
		&r.NoTelp,
		&r.NoFax,
		&r.Email,
		&r.Website,
		&r.Narahubung,
		&r.NPWP,
		&r.NPWPPemberi,
		&r.NPWPFile,
		&r.SPIFile,
		&r.Logo,
		&r.CreatedAt,
		&r.CreatedBy,
		&r.UpdatedAt,
		&r.UpdatedBy,
	)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data BUMD. - " + err.Error(),
		}
	}

	q = `SELECT nama_daerah, id_prop FROM data.m_daerah WHERE id_daerah = $1`
	err = c.pgxConnMstData.QueryRow(fCtx, q, r.IdDaerah).Scan(&r.NamaDaerah, &r.IdProvinsi)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Daerah. - " + err.Error(),
		}
	}

	if r.IdDaerah != r.IdProvinsi {
		q = `SELECT nama_daerah FROM data.m_daerah WHERE id_daerah = $1`
		err = c.pgxConnMstData.QueryRow(fCtx, q, r.IdProvinsi).Scan(&r.NamaProvinsi)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data Daerah. - " + err.Error(),
			}
		}
	}

	return r, err
}

func (c *BumdController) Create(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.BumdForm,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat BUMD. - " + err.Error(),
		}
	}

	q := `
	INSERT INTO trn_bumd (
		id_bumd,
		id_daerah,
		id_bentuk_hukum,
		id_bentuk_usaha,
		id_induk_perusahaan,
		penerapan_spi_bumd,
		nama_bumd,
		deskripsi_bumd,
		alamat_bumd,
		no_telp_bumd,
		no_fax_bumd,
		email_bumd,
		website_bumd,
		narahubung_bumd,
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
		$15
	)
	`

	_, err = c.pgxConn.Exec(
		fCtx,
		q,
		id,                        // 1
		idDaerah,                  // 2
		payload.IdBentukHukum,     // 3
		payload.IdBentukUsaha,     // 4
		payload.IdIndukPerusahaan, // 5
		payload.PenerapanSPI,      // 6
		payload.Nama,              // 7
		payload.Deskripsi,         // 8
		payload.Alamat,            // 9
		payload.NoTelp,            // 10
		payload.NoFax,             // 11
		payload.Email,             // 12
		payload.Website,           // 13
		payload.Narahubung,        // 14
		idUser,                    // 15
	)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat BUMD. - " + err.Error(),
		}
	}

	return true, err
}

func (c *BumdController) Update(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.BumdForm,
	id uuid.UUID,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM trn_bumd WHERE id_bumd = $1 AND deleted_by = 0
	`
	args = append(args, id)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	var count int
	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengubah BUMD. - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusNotFound,
			Message: "Data BUMD tidak ditemukan.",
		}
	}

	args = []interface{}{}
	q = `
	UPDATE trn_bumd
	SET
		id_bentuk_hukum = $1,
		id_bentuk_usaha = $2,
		id_induk_perusahaan = $3,
		penerapan_spi_bumd = $4,
		nama_bumd = $5,
		deskripsi_bumd = $6,
		alamat_bumd = $7,
		no_telp_bumd = $8,
		no_fax_bumd = $9,
		email_bumd = $10,
		website_bumd = $11,
		narahubung_bumd = $12,
		updated_by = $13,
		updated_at = NOW()
	WHERE id_bumd = $14 AND deleted_by = 0
	`
	args = append(args,
		payload.IdBentukHukum,     // 1
		payload.IdBentukUsaha,     // 2
		payload.IdIndukPerusahaan, // 3
		payload.PenerapanSPI,      // 4
		payload.Nama,              // 5
		payload.Deskripsi,         // 6
		payload.Alamat,            // 7
		payload.NoTelp,            // 8
		payload.NoFax,             // 9
		payload.Email,             // 10
		payload.Website,           // 11
		payload.Narahubung,        // 12
		idUser,                    // 13
		id,                        // 14
	)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	_, err = c.pgxConn.Exec(
		fCtx,
		q,
		args...,
	)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus BUMD - " + err.Error(),
		}
	}

	return true, err
}

func (c *BumdController) Delete(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	id uuid.UUID,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM trn_bumd WHERE id_bumd = $1 AND deleted_by = 0
	`
	args = append(args, id)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	var count int
	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus BUMD - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusNotFound,
			Message: "Data BUMD tidak ditemukan.",
		}
	}

	q = `
	SELECT COALESCE(COUNT(*), 0) FROM trn_bumd WHERE id_induk_perusahaan = $1 AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus BUMD - " + err.Error(),
		}
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Data BUMD induk tidak dapat dihapus, karena masih ada data anak BUMD yang terkait.",
		}
	}

	q = `
	UPDATE trn_bumd
		SET deleted_by = $1, deleted_at = NOW()
	WHERE id_bumd = $2
	`
	args = append(args, idUser, id)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	_, err = c.pgxConn.Exec(fCtx, q, args...)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus BUMD - " + err.Error(),
		}
	}

	return true, err
}

func (c *BumdController) Logo(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	id uuid.UUID,
) (r bumd.LogoModel, err error) {
	q := `
	SELECT logo_bumd FROM trn_bumd WHERE id_bumd = $1 AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.Logo)
	if err != nil {
		return bumd.LogoModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengambil data logo - " + err.Error(),
		}
	}
	return r, nil
}

func (c *BumdController) LogoUpdate(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.LogoForm,
	id uuid.UUID,
) (r bool, err error) {
	if payload.Logo != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.Logo.Filename)

		src, err := payload.Logo.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
		}
		defer src.Close()

		// upload file
		objectName := "bumd_logo/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.Logo.Size,
			minio.PutObjectOptions{ContentType: payload.Logo.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q := `UPDATE trn_bumd SET logo_bumd=$1 WHERE id_bumd=$2`
		_, err = c.pgxConn.Exec(fCtx, q, objectName, id)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *BumdController) SPI(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	id uuid.UUID,
) (r bumd.SPIModel, err error) {
	q := `
		SELECT penerapan_spi_bumd, file_spi_bumd FROM trn_bumd WHERE id_bumd = $1 AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.PenerapanSPI, &r.FileSPI)
	if err != nil {
		return bumd.SPIModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data SPI. - " + err.Error(),
		}
	}
	return r, nil
}

func (c *BumdController) SPIUpdate(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.SPIForm,
	id uuid.UUID,
) (r bool, err error) {
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
	UPDATE trn_bumd
	SET
		penerapan_spi_bumd = $1
	WHERE id_bumd = $2
	`
	_, err = tx.Exec(context.Background(), q, payload.PenerapanSPI, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate SPI. - " + err.Error(),
		}
	}

	if payload.PenerapanSPI && payload.FileSPI != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.FileSPI.Filename)

		src, err := payload.FileSPI.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
		}
		defer src.Close()

		// upload file
		objectName := "bumd_spi/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.FileSPI.Size,
			minio.PutObjectOptions{ContentType: payload.FileSPI.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		// update file
		q = `UPDATE trn_bumd SET file_spi_bumd=$1 WHERE id_bumd=$2`
		_, err = tx.Exec(context.Background(), q, objectName, id)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}
	return true, err
}

func (c *BumdController) NPWP(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	id uuid.UUID,
) (r bumd.NPWPModel, err error) {
	q := `
	SELECT npwp_bumd, npwp_pemberi_bumd, npwp_file_bumd FROM trn_bumd WHERE id_bumd = $1 AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.NPWP, &r.Pemberi, &r.File)
	if err != nil {
		return bumd.NPWPModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengambil data NPWP - " + err.Error(),
		}
	}
	return r, nil
}

func (c *BumdController) NPWPUpdate(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.NPWPForm,
	id uuid.UUID,
) (r bool, err error) {
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
	UPDATE trn_bumd
	SET
		npwp_bumd = $1,
		npwp_pemberi_bumd = $2
	WHERE id_bumd = $3
	`
	_, err = tx.Exec(context.Background(), q, payload.NPWP, payload.Pemberi, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate NPWP. - " + err.Error(),
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
		objectName := "bumd_npwp/" + fileName
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
		q = `UPDATE trn_bumd SET npwp_file_bumd=$1 WHERE id_bumd=$2`
		_, err = tx.Exec(context.Background(), q, objectName, id)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *BumdController) KelengkapanInput(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idDaerah int,
	bentukUsaha,
	idBumd uuid.UUID,
	page,
	limit int,
	search string,
) (r []bumd.KelengkapanInputModel, totalCount, pageCount int, err error) {
	r = make([]bumd.KelengkapanInputModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idDaerahClaims := int(claims["id_daerah"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Kelengkapan Input. - " + err.Error(),
		}
	}

	offset := limit * (page - 1)

	qCount := `
	SELECT COALESCE(COUNT(*), 0) FROM trn_bumd WHERE deleted_by = 0
	`
	args := make([]interface{}, 0)

	if idDaerah > 0 {
		qCount += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
	}
	if idBumd != uuid.Nil {
		qCount += fmt.Sprintf(` AND id_bumd = $%d`, len(args)+1)
	}
	if search != "" {
		qCount += fmt.Sprintf(` AND nama_bumd ILIKE $%d`, len(args)+1)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data Kelengkapan Input. - " + err.Error(),
		}
	}

	q := `
	WITH m_daerah_temp AS (
		SELECT * FROM dblink ($1,
		'
			SELECT id_daerah, nama_daerah
			FROM data.m_daerah
			WHERE is_deleted = 0
		')
	AS m_daerah_temp
		(id_daerah INT4, nama_daerah VARCHAR)
	)
	SELECT
		b.id_bumd,
		b.nama_bumd,
		b.id_bentuk_hukum,
		mbbh.nama_bbh as bentuk_badan_hukum,
		b.id_bentuk_usaha,
		mbbu.nama_bu as bentuk_usaha,
		mda.nama_daerah as nama_daerah,
		mda.id_daerah as id_daerah,
		CASE WHEN
			penerapan_spi_bumd = true
			AND file_spi_bumd != ''
			THEN 1
			ELSE 0
		END as penerapan_spi,
		CASE WHEN
			(SELECT COUNT(*) FROM trn_perda_pendirian WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as akta_pendirian,
		CASE WHEN
			(SELECT COUNT(*) FROM trn_kinerja WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as kinerja,
		1 as keuangan,
		/* CASE WHEN
			(SELECT COUNT(*) FROM trn_keuangan WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as keuangan, */
		CASE WHEN
			(SELECT COUNT(*) FROM trn_modal WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as modal,
		CASE WHEN
			(SELECT COUNT(*) FROM trn_pegawai WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as pegawai,
		CASE WHEN
			(SELECT COUNT(*) FROM trn_pengurus WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as pengurus,
		CASE WHEN
			(SELECT COUNT(*) FROM trn_peraturan WHERE id_bumd = b.id_bumd AND deleted_by = 0) > 0
			THEN 1
			ELSE 0
		END as peraturan
	FROM trn_bumd b
		LEFT JOIN m_bentuk_badan_hukum mbbh ON mbbh.id_bbh = b.id_bentuk_hukum
		LEFT JOIN m_bentuk_usaha mbbu ON mbbu.id_bu = b.id_bentuk_usaha
		LEFT JOIN m_daerah_temp mda ON mda.id_daerah = b.id_daerah
	WHERE b.deleted_by = 0
	`
	args = make([]interface{}, 0)
	args = append(args, os.Getenv("DB_SERVER_URL_MST_DATA"))

	if idDaerahClaims > 0 {
		idDaerah = idDaerahClaims
	}
	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}
	if idBumd != uuid.Nil {
		q += fmt.Sprintf(` AND id_bumd = $%d`, len(args)+1)
		args = append(args, idBumd)
	}

	if search != "" {
		q += fmt.Sprintf(` AND nama_bumd ILIKE $%d`, len(args)+1)
		args = append(args, "%"+search+"%")
	}

	q += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Kelengkapan Input. - " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		var m bumd.KelengkapanInputModel
		err = rows.Scan(
			&m.IdBumd,
			&m.NamaBumd,
			&m.IdBentukBadanHukum,
			&m.BentukBadanHukum,
			&m.IdBentukUsaha,
			&m.BentukUsaha,
			&m.NamaDaerah,
			&m.IdDaerah,
			&m.PenerapanSPI,
			&m.AktaPendirian,
			&m.Kinerja,
			&m.Keuangan,
			&m.Modal,
			&m.Pegawai,
			&m.Pengurus,
			&m.Peraturan,
		)

		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data Kelengkapan Input. - " + err.Error(),
			}
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, nil
}

func (c *BumdController) Sebaran(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idDaerah int,
) (r bumd.SebaranModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idDaerahClaims := int(claims["id_daerah"].(float64))
	if idDaerah > 0 {
		idDaerah = idDaerahClaims
	}

	r.IdDaerah = int32(idDaerah)

	q := `
	SELECT
		sum(case when id_bentuk_usaha = '01994c01-c285-7eea-aead-221a0d1f4cac' then 1 else 0 end) as bpd,
		sum(case when id_bentuk_usaha = '01994c01-c285-7e6f-865a-b286738dff03' then 1 else 0 end) as bpr,
		sum(case when id_bentuk_usaha = '01994c01-c285-73ae-8885-65ab6b20deb3' then 1 else 0 end) as jamkrida,
		sum(case when id_bentuk_usaha = '01994c01-c285-7e57-a486-fd9978083917' then 1 else 0 end) as pdam,
		sum(case when id_bentuk_usaha = '01994c01-c285-7e51-9140-e75269630f28' then 1 else 0 end) as pasar,
		sum(case when id_bentuk_usaha = '01994c01-c285-7744-ba6a-5782f39fe366' then 1 else 0 end) as aneka_usaha,
		sum(case when id_bentuk_usaha = '01994c01-c285-73b3-ada9-10d5180a4a2f' then 1 else 0 end) as lainnya
	FROM trn_bumd
	WHERE deleted_by = 0
	`

	var args []interface{}
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&r.Bpd, &r.Bpr, &r.Jamkrida, &r.Pdam, &r.Pasar, &r.AnekaUsaha, &r.Lainnya)
	if err != nil {
		return bumd.SebaranModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Sebaran BUMD. - " + err.Error(),
		}
	}
	return r, nil
}
