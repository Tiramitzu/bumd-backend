package controller

import (
	"fmt"
	"math"
	models "microdata/kemendagri/bumd/model"
	"microdata/kemendagri/bumd/utils"
	"net/http"
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
) (r []models.BumdModel, totalCount, pageCount int, err error) {
	r = make([]models.BumdModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id_daerah = $1 AND deleted_by = 0
	`
	args = append(args, idDaerah)
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data BUMD: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT
		b.id,
		b.id_daerah,
		b.id_bentuk_hukum,
		b.id_bentuk_usaha,
		mbbh.nama as bentuk_badan_hukum,
		mbbu.nama as bentuk_usaha,
		b.nama,
		b.deskripsi,
		b.alamat,
		b.no_telp,
		b.email,
		b.website
	FROM bumd b
		LEFT JOIN mst_bentuk_badan_hukum mbbh ON mbbh.id = b.id_bentuk_hukum
		LEFT JOIN mst_bentuk_usaha mbbu ON mbbu.id = b.id_bentuk_usaha
	WHERE b.id_daerah = $1
		AND b.deleted_by = 0
	`
	args = append(args, idDaerah)
	if nama != "" {
		q += fmt.Sprintf(` AND b.nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
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
		var m models.BumdModel
		err = rows.Scan(
			&m.ID,
			&m.IDDaerah,
			&m.IDBentukHukum,
			&m.IDBentukUsaha,
			&m.BentukBadanHukum,
			&m.BentukUsaha,
			&m.Nama,
			&m.Deskripsi,
			&m.Alamat,
			&m.NoTelp,
			&m.Email,
			&m.Website,
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

func (c *BumdController) Create(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *models.BumdForm,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	q := `
	INSERT INTO bumd (
		id_daerah,
		id_bentuk_hukum,
		id_bentuk_usaha,
		nama,
		deskripsi,
		alamat,
		no_telp,
		email,
		website,
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
		$10
	)
	`

	_, err = c.pgxConn.Exec(
		fCtx,
		q,
		idDaerah,              // 1
		payload.IDBentukHukum, // 2
		payload.IDBentukUsaha, // 3
		payload.Nama,          // 4
		payload.Deskripsi,     // 5
		payload.Alamat,        // 6
		payload.NoTelp,        // 7
		payload.Email,         // 8
		payload.Website,       // 9
		idUser,                // 10
	)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *BumdController) Update(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	payload *models.BumdForm,
	id int,
) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id = $1 AND id_daerah = $2
	`

	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id, idDaerah).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, utils.RequestError{
			Code:    http.StatusNotFound,
			Message: "Data BUMD tidak ditemukan.",
		}
	}

	q = `
	UPDATE bumd
	SET
		id_bentuk_hukum = $1,
		id_bentuk_usaha = $2,
		nama = $3,
		deskripsi = $4,
		alamat = $5,
		no_telp = $6,
		email = $7,
		website = $8,
		updated_by = $9
	WHERE id = $10 AND id_daerah = $11
	`

	_, err = c.pgxConn.Exec(
		fCtx,
		q,
		payload.IDBentukHukum,
		payload.IDBentukUsaha,
		payload.Nama,
		payload.Deskripsi,
		payload.Alamat,
		payload.NoTelp,
		payload.Email,
		payload.Website,
		idUser,
		id,
		idDaerah,
	)
	if err != nil {
		return false, err
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

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM bumd WHERE id = $1 AND id_daerah = $2
	`

	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id, idDaerah).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, utils.RequestError{
			Code:    http.StatusNotFound,
			Message: "Data BUMD tidak ditemukan.",
		}
	}

	q = `
	UPDATE bumd
		SET deleted_by = $1, deleted_at = $2
	WHERE id = $3 AND id_daerah = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id, idDaerah)
	if err != nil {
		return false, err
	}

	return true, err
}
