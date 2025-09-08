package bumd

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type BumdController struct {
	pgxConn *pgxpool.Pool
}

func NewBumdController(pgxConn *pgxpool.Pool) *BumdController {
	return &BumdController{pgxConn: pgxConn}
}

func (c *BumdController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	page, limit int,
	nama string,
	penerapanSPI bool,
	indukPerusahaan int,
) (r []bumd.BumdModel, totalCount, pageCount int, err error) {
	r = make([]bumd.BumdModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))
	offset := limit * (page - 1)

	var args []interface{}
	qCount := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE deleted_by = 0
	`

	q := `
	WITH t_induk_perusahaan AS (
		SELECT id, nama as nama_induk_perusahaan
		FROM bumd
		WHERE deleted_by = 0
	)
	SELECT
		b.id,
		b.id_daerah,
		b.id_bentuk_hukum,
		b.id_bentuk_usaha,
		b.id_induk_perusahaan,
		COALESCE(t_induk_perusahaan.nama_induk_perusahaan, '-') as nama_induk_perusahaan,
		b.penerapan_spi,
		mbbh.nama as bentuk_badan_hukum,
		mbbu.nama as bentuk_usaha,
		b.nama,
		b.deskripsi,
		b.alamat,
		b.no_telp,
		b.no_fax,
		b.email,
		b.website,
		b.narahubung
	FROM bumd b
		LEFT JOIN mst_bentuk_badan_hukum mbbh ON mbbh.id = b.id_bentuk_hukum
		LEFT JOIN mst_bentuk_usaha mbbu ON mbbu.id = b.id_bentuk_usaha
		LEFT JOIN t_induk_perusahaan ON t_induk_perusahaan.id = b.id_induk_perusahaan
	WHERE b.deleted_by = 0
	`
	if idDaerah > 0 {
		qCount += fmt.Sprintf(` ADN b.id_daerah = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}
	if nama != "" {
		qCount += fmt.Sprintf(` AND b.nama ILIKE $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}
	if penerapanSPI {
		qCount += fmt.Sprintf(` AND b.penerapan_spi = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.penerapan_spi = $%d`, len(args)+1)
		args = append(args, penerapanSPI)
	}
	if indukPerusahaan != 0 {
		qCount += fmt.Sprintf(` AND b.id_induk_perusahaan = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND b.id_induk_perusahaan = $%d`, len(args)+1)
		args = append(args, indukPerusahaan)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data BUMD: %w", err)
	}

	q += fmt.Sprintf(`
	ORDER BY id DESC
	LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data BUMD: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m bumd.BumdModel
		err = rows.Scan(
			&m.ID,
			&m.IDDaerah,
			&m.IDBentukHukum,
			&m.IDBentukUsaha,
			&m.IDIndukPerusahaan,
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
		)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data BUMD: %w", err)
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
	id int,
) (r bumd.BumdModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	WITH t_induk_perusahaan AS (
		SELECT id, nama as nama_induk_perusahaan
		FROM bumd
		WHERE deleted_by = 0
	)
	SELECT
		b.id,
		b.id_daerah,
		b.id_bentuk_hukum,
		b.id_bentuk_usaha,
		b.id_induk_perusahaan,
		COALESCE(t_induk_perusahaan.nama_induk_perusahaan, '-') as nama_induk_perusahaan,
		b.penerapan_spi,
		mbbh.nama as bentuk_badan_hukum,
		mbbu.nama as bentuk_usaha,
		b.nama,
		b.deskripsi,
		b.alamat,
		b.no_telp,
		b.no_fax,
		b.email,
		b.website,
		b.narahubung
	FROM bumd b
		LEFT JOIN mst_bentuk_badan_hukum mbbh ON mbbh.id = b.id_bentuk_hukum
		LEFT JOIN mst_bentuk_usaha mbbu ON mbbu.id = b.id_bentuk_usaha
		LEFT JOIN t_induk_perusahaan ON t_induk_perusahaan.id = b.id_induk_perusahaan
	WHERE b.id = $1 AND b.deleted_by = 0
	`
	args = append(args, id)
	if idDaerah > 0 {
		q += fmt.Sprintf(` AND b.id_daerah = $%d`, len(args)+1)
		args = append(args, idDaerah)
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(
		&r.ID,
		&r.IDDaerah,
		&r.IDBentukHukum,
		&r.IDBentukUsaha,
		&r.IDIndukPerusahaan,
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
	)
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data BUMD: %w", err)
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

	q := `
	INSERT INTO bumd (
		id_daerah,
		id_bentuk_hukum,
		id_bentuk_usaha,
		id_induk_perusahaan,
		penerapan_spi,
		nama,
		deskripsi,
		alamat,
		no_telp,
		no_fax,
		email,
		website,
		narahubung,
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
		$14
	)
	`

	_, err = c.pgxConn.Exec(
		fCtx,
		q,
		idDaerah,                  // 1
		payload.IDBentukHukum,     // 2
		payload.IDBentukUsaha,     // 3
		payload.IDIndukPerusahaan, // 4
		payload.PenerapanSPI,      // 5
		payload.Nama,              // 6
		payload.Deskripsi,         // 7
		payload.Alamat,            // 8
		payload.NoTelp,            // 9
		payload.NoFax,             // 10
		payload.Email,             // 11
		payload.Website,           // 12
		payload.Narahubung,        // 13
		idUser,                    // 14
	)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat BUMD - " + err.Error(),
		}
	}

	return true, err
}

func (c *BumdController) Update(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *bumd.BumdForm,
	id int,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id = $1 AND deleted_by = 0
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
			Message: "Gagal mengubah BUMD - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusNotFound,
			Message: "Data BUMD tidak ditemukan.",
		}
	}

	q = `
	UPDATE bumd
	SET
		id_bentuk_hukum = $1,
		id_bentuk_usaha = $2,
		id_induk_perusahaan = $3,
		penerapan_spi = $4,
		nama = $5,
		deskripsi = $6,
		alamat = $7,
		no_telp = $8,
		no_fax = $9,
		email = $10,
		website = $11,
		narahubung = $12,
		updated_by = $13
	WHERE id = $14 AND deleted_by = 0
	`
	args = append(args,
		payload.IDBentukHukum,     // 1
		payload.IDBentukUsaha,     // 2
		payload.IDIndukPerusahaan, // 3
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
	id int,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id = $1 AND deleted_by = 0
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
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id_induk_perusahaan = $1 AND deleted_by = 0
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
			Message: "Data BUMD induk tidak dapat dihapus, karena masih ada data BUMD anak yang terkait.",
		}
	}

	q = `
	UPDATE bumd
		SET deleted_by = $1, deleted_at = $2
	WHERE id = $3 AND deleted_by = 0
	`
	args = append(args, idUser, time.Now(), id)
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
