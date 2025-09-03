package controller_mst

import (
	"fmt"
	"math"
	models "microdata/kemendagri/bumd/model/master"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type BentukUsahaController struct {
	pgxConn *pgxpool.Pool
}

func NewBentukUsahaController(pgxConn *pgxpool.Pool) *BentukUsahaController {
	return &BentukUsahaController{pgxConn: pgxConn}
}

func (c *BentukUsahaController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, nama string) (r []models.BentukUsahaModel, totalCount, pageCount int, err error) {
	r = make([]models.BentukUsahaModel, 0)
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_usaha WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Bentuk Usaha: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, nama, deskripsi FROM mst_bentuk_usaha WHERE deleted_by = 0
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
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m models.BentukUsahaModel
		err = rows.Scan(&m.ID, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Bentuk Usaha: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_usaha WHERE nama = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha sudah ada.",
		}
	}

	q = `
	INSERT INTO mst_bentuk_usaha (nama, deskripsi, created_by) VALUES ($1, $2, $3)
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *BentukUsahaController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukUsahaForm, id int) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_usaha WHERE nama = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_bentuk_usaha
	SET nama = $1, deskripsi = $2, updated_by = $3
	WHERE id = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *BentukUsahaController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_usaha WHERE id = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Usaha tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_bentuk_usaha
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, err
	}

	return true, err
}
