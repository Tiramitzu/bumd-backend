package controller_mst

import (
	"fmt"
	"math"
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type RolesController struct {
	pgxConn *pgxpool.Pool
}

func NewRolesController(pgxConn *pgxpool.Pool) *RolesController {
	return &RolesController{pgxConn: pgxConn}
}

func (c *RolesController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, nama string) (r []models.RolesModel, totalCount, pageCount int, err error) {
	r = make([]models.RolesModel, 0)
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM roles WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data Roles: " + err.Error(),
		}
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, nama, deskripsi FROM roles WHERE deleted_by = 0
	`
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}
	q += fmt.Sprintf(`
	LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Roles: " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m models.RolesModel
		err = rows.Scan(&m.ID, &m.Nama)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data Roles: " + err.Error(),
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

func (c *RolesController) View(fCtx *fasthttp.RequestCtx, id int) (r models.RolesModel, err error) {
	q := `
	SELECT id, nama, deskripsi
	FROM roles
	WHERE id = $1
	AND deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.ID, &r.Nama)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Roles: " + err.Error(),
		}
	}

	return r, err
}
