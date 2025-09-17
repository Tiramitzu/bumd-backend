package controller_mst

import (
	"fmt"
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type JenisDokumenController struct {
	pgxConn *pgxpool.Pool
}

func NewJenisDokumenController(pgxConn *pgxpool.Pool) *JenisDokumenController {
	return &JenisDokumenController{pgxConn: pgxConn}
}

func (c *JenisDokumenController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	nama string,
) (
	r []models.JenisDokumenModel,
	err error,
) {
	r = make([]models.JenisDokumenModel, 0)

	var args []interface{}
	q := `
	SELECT id_jd, nama_jd, deskripsi_jd FROM m_jenis_dokumen WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama_jd ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	q += `ORDER BY created_at DESC`

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "data Jenis Dokumen tidak ditemukan",
			}
		}
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Jenis Dokumen: " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m models.JenisDokumenModel
		err = rows.Scan(&m.Id, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data Jenis Dokumen: " + err.Error(),
			}
		}
		r = append(r, m)
	}

	return r, err
}

func (c *JenisDokumenController) View(fCtx *fasthttp.RequestCtx, id uuid.UUID) (r models.JenisDokumenModel, err error) {
	q := `
	SELECT id_jd, nama_jd, deskripsi_jd
	FROM m_jenis_dokumen 
	WHERE id_jd = $1 AND deleted_by = 0
	AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.Id, &r.Nama, &r.Deskripsi)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "data Jenis Dokumen tidak ditemukan",
			}
		}
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Jenis Dokumen: " + err.Error(),
		}
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

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat data Jenis Dokumen - " + err.Error(),
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_jenis_dokumen WHERE nama_jd = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen sudah ada.",
		}
	}

	q = `
	INSERT INTO m_jenis_dokumen (id_jd, nama_jd, deskripsi_jd, created_by) VALUES ($1, $2, $3, $4)
	`

	_, err = c.pgxConn.Exec(fCtx, q, id, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}

func (c *JenisDokumenController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.JenisDokumenForm, id uuid.UUID) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM m_jenis_dokumen WHERE nama_jd = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_jenis_dokumen
	SET nama_jd = $1, deskripsi_jd = $2, updated_by = $3
	WHERE id_jd = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memperbarui data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}

func (c *JenisDokumenController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id uuid.UUID) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM m_jenis_dokumen WHERE id_jd = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung data Jenis Dokumen - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Jenis Dokumen tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_jenis_dokumen
	SET deleted_by = $1, deleted_at = $2
	WHERE id_jd = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghapus data Jenis Dokumen - " + err.Error(),
		}
	}

	return true, err
}
