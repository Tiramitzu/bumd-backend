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

type BentukUsahaController struct {
	pgxConn *pgxpool.Pool
}

func NewBentukUsahaController(pgxConn *pgxpool.Pool) *BentukUsahaController {
	return &BentukUsahaController{pgxConn: pgxConn}
}

func (c *BentukUsahaController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, nama string) (r []models.BentukUsahaModel, err error) {
	r = make([]models.BentukUsahaModel, 0)

	var args []interface{}
	q := `
	SELECT id_bu, nama_bu, deskripsi_bu FROM m_bentuk_usaha WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama_bu ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	q += ` ORDER BY created_at DESC`

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, fmt.Errorf("data Bentuk Usaha tidak ditemukan")
		}
		return r, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m models.BentukUsahaModel
		err = rows.Scan(&m.Id, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, fmt.Errorf("gagal memindahkan data Bentuk Usaha: %w", err)
		}
		r = append(r, m)
	}

	return r, err
}

func (c *BentukUsahaController) View(fCtx *fasthttp.RequestCtx, id uuid.UUID) (r models.BentukUsahaModel, err error) {
	q := `
	SELECT id_bu, nama_bu, deskripsi_bu
	FROM m_bentuk_usaha
	WHERE id_bu = $1
	AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.Id, &r.Nama, &r.Deskripsi)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, fmt.Errorf("data Bentuk Usaha tidak ditemukan")
		}
		return r, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	return r, err
}

func (c *BentukUsahaController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukUsahaForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk membuat data Bentuk Usaha.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_usaha WHERE nama_bu = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Bentuk Usaha - " + err.Error(),
		}
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha sudah ada.",
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat data Bentuk Usaha - " + err.Error(),
		}
	}

	q = `
	INSERT INTO m_bentuk_usaha (id_bu, nama_bu, deskripsi_bu, created_by) VALUES ($1, $2, $3, $4)
	`

	_, err = c.pgxConn.Exec(fCtx, q, id, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat data Bentuk Usaha - " + err.Error(),
		}
	}

	return true, err
}

func (c *BentukUsahaController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukUsahaForm, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk memperbarui data Bentuk Usaha.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_usaha WHERE nama_bu = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data Bentuk Usaha - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_bentuk_usaha
	SET nama_bu = $1, deskripsi_bu = $2, updated_by = $3
	WHERE id_bu = $4 AND deleted_by = 0
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal memperbarui data Bentuk Usaha - " + err.Error(),
		}
	}

	return true, err
}

func (c *BentukUsahaController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk menghapus data Bentuk Usaha.",
		}
	}

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM m_bentuk_usaha WHERE id_bu = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data Bentuk Usaha - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha tidak ditemukan.",
		}
	}

	q = `
	UPDATE m_bentuk_usaha
	SET deleted_by = $1, deleted_at = $2
	WHERE id_bu = $3 AND deleted_by = 0
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data Bentuk Usaha - " + err.Error(),
		}
	}

	return true, err
}
