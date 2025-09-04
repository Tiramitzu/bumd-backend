package controller_mst

import (
	"fmt"
	"math"
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type JenisDokumenController struct {
	pgxConn *pgxpool.Pool
}

func NewJenisDokumenController(pgxConn *pgxpool.Pool) *JenisDokumenController {
	return &JenisDokumenController{pgxConn: pgxConn}
}

func (c *JenisDokumenController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, nama string) (r []models.JenisDokumenModel, totalCount, pageCount int, err error) {
	r = make([]models.JenisDokumenModel, 0)
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_jenis_dokumen WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Jenis Dokumen: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, nama, deskripsi FROM mst_jenis_dokumen WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}
	q += fmt.Sprintf(`
	ORDER BY id DESC
	LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, totalCount, pageCount, fmt.Errorf("data Jenis Dokumen tidak ditemukan")
		}
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Jenis Dokumen: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m models.JenisDokumenModel
		err = rows.Scan(&m.ID, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Jenis Dokumen: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *JenisDokumenController) View(fCtx *fasthttp.RequestCtx, id int) (r models.JenisDokumenModel, err error) {
	q := `
	SELECT id, nama, deskripsi
	FROM mst_jenis_dokumen 
	WHERE id = $1
	AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.ID, &r.Nama, &r.Deskripsi)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, fmt.Errorf("data Jenis Dokumen tidak ditemukan")
		}
		return r, fmt.Errorf("gagal mengambil data Jenis Dokumen: %w", err)
	}

	return r, err
}

func (c *JenisDokumenController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.JenisDokumenForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk membuat data Jenis Dokumen.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_jenis_dokumen WHERE nama = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen sudah ada.",
		}
	}

	q = `
	INSERT INTO mst_jenis_dokumen (nama, deskripsi, created_by) VALUES ($1, $2, $3)
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}

func (c *JenisDokumenController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.JenisDokumenForm, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk memperbarui data Jenis Dokumen.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_jenis_dokumen WHERE nama = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_jenis_dokumen
	SET nama = $1, deskripsi = $2, updated_by = $3
	WHERE id = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal memperbarui data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}

func (c *JenisDokumenController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk menghapus data Jenis Dokumen.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_jenis_dokumen WHERE id = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_jenis_dokumen
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}
