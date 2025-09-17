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

type BentukBadanHukumController struct {
	pgxConn *pgxpool.Pool
}

func NewBentukBadanHukumController(pgxConn *pgxpool.Pool) *BentukBadanHukumController {
	return &BentukBadanHukumController{pgxConn: pgxConn}
}

func (c *BentukBadanHukumController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, nama string) (r []models.BentukBadanHukumModel, err error) {
	r = make([]models.BentukBadanHukumModel, 0)

	var args []interface{}
	q := `
	SELECT id_bbh, nama_bbh, deskripsi_bbh FROM m_bentuk_badan_hukum WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama_bbh ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	q += ` ORDER BY created_at DESC`

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "data Bentuk Badan Hukum tidak ditemukan",
			}
		}
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Bentuk Badan Hukum: " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m models.BentukBadanHukumModel
		err = rows.Scan(&m.Id, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data Bentuk Badan Hukum: " + err.Error(),
			}
		}
		r = append(r, m)
	}

	return r, err
}

func (c *BentukBadanHukumController) View(fCtx *fasthttp.RequestCtx, id uuid.UUID) (r models.BentukBadanHukumModel, err error) {
	q := `
	SELECT id_bbh, nama_bbh, deskripsi_bbh
	FROM m_bentuk_badan_hukum
	WHERE id_bbh = $1
	AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.Id, &r.Nama, &r.Deskripsi)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "data Bentuk Badan Hukum tidak ditemukan",
			}
		}
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Bentuk Badan Hukum: " + err.Error(),
		}
	}

	return r, err
}

func (c *BentukBadanHukumController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukBadanHukumForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk membuat data Bentuk Badan Hukum.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_badan_hukum WHERE nama_bbh = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum sudah ada.",
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	q = `
	INSERT INTO m_bentuk_badan_hukum (id_bbh, nama_bbh, deskripsi_bbh, created_by) VALUES ($1, $2, $3, $4)
	`

	_, err = c.pgxConn.Exec(fCtx, q, id, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	return true, err
}

func (c *BentukBadanHukumController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukBadanHukumForm, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk memperbarui data Bentuk Badan Hukum.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_badan_hukum WHERE nama_bbh = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_bentuk_badan_hukum
	SET nama_bbh = $1, deskripsi_bbh = $2, updated_by = $3
	WHERE id_bbh = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal memperbarui data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	return true, err
}

func (c *BentukBadanHukumController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk menghapus data Bentuk Badan Hukum.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_badan_hukum WHERE id_bbh = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_bentuk_badan_hukum
	SET deleted_by = $1, deleted_at = $2
	WHERE id_bbh = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data Bentuk Badan Hukum - " + err.Error(),
		}
	}

	return true, err
}
