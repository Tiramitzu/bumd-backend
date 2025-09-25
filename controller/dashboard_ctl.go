package controller

import (
	"microdata/kemendagri/bumd/models"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type DashboardController struct {
	pgxConn *pgxpool.Pool
}

func NewDashboardController(pgxConn *pgxpool.Pool) *DashboardController {
	return &DashboardController{pgxConn: pgxConn}
}

func (c *DashboardController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token) (r models.DashboardModel, err error) {
	q := `
	SELECT COUNT(*) AS total_bumd FROM trn_bumd WHERE deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q).Scan(&r.TotalBumd)
	if err != nil {
		return models.DashboardModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Dashboard. - " + err.Error(),
		}
	}

	q = `
	SELECT SUM(jumlah_modal) AS total_modal FROM trn_modal WHERE deleted_by = 0
	`
	err = c.pgxConn.QueryRow(fCtx, q).Scan(&r.TotalModal)
	if err != nil {
		return models.DashboardModel{}, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Dashboard. - " + err.Error(),
		}
	}

	return r, nil
}
