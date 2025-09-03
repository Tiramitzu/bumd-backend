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

type BentukBadanHukumController struct {
	pgxConn *pgxpool.Pool
}

func NewBentukBadanHukumController(pgxConn *pgxpool.Pool) *BentukBadanHukumController {
	return &BentukBadanHukumController{pgxConn: pgxConn}
}

func (c *BentukBadanHukumController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, nama string) (r []models.BentukBadanHukumModel, totalCount, pageCount int, err error) {
	r = make([]models.BentukBadanHukumModel, 0)

	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_badan_hukum WHERE deleted_by = 0
	`
	args = append(args, idDaerah)
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Bentuk Badan Hukum: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, nama, deskripsi FROM mst_bentuk_badan_hukum WHERE deleted_by = 0
	`
	args = append(args, idDaerah)
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
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Badan Hukum: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m models.BentukBadanHukumModel
		err = rows.Scan(&m.ID, &m.Nama, &m.Deskripsi)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Bentuk Badan Hukum: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_badan_hukum WHERE nama = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum sudah ada.",
		}
	}

	q = `
	INSERT INTO mst_bentuk_badan_hukum (nama, deskripsi, created_by) VALUES ($1, $2, $3)
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *BentukBadanHukumController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.BentukBadanHukumForm, id int) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_badan_hukum WHERE nama = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, payload.Nama).Scan(&count)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_bentuk_badan_hukum
	SET nama = $1, deskripsi = $2, updated_by = $3
	WHERE id = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nama, payload.Deskripsi, idUser, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *BentukBadanHukumController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r bool, err error) {
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
	SELECT COALESCE(COUNT(*), 0) FROM mst_bentuk_badan_hukum WHERE id = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "Bentuk Badan Hukum tidak ditemukan.",
		}
	}

	q = `
	UPDATE mst_bentuk_badan_hukum
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, err
	}

	return true, err
}
