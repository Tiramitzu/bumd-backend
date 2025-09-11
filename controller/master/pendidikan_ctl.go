package controller_mst

import (
	"context"
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type PendidikanController struct {
	contextTimeout time.Duration
	pgxConn        *pgxpool.Pool
}

func NewPendidikanController(conn *pgxpool.Pool, tot time.Duration) *PendidikanController {
	return &PendidikanController{
		pgxConn:        conn,
		contextTimeout: tot,
	}
}

func (c *PendidikanController) Index(ctx context.Context, jwt *jwt.Token) ([]models.PendidikanModel, error) {
	var r []models.PendidikanModel

	q := `SELECT id, nama FROM m_pendidikan`
	rows, err := c.pgxConn.Query(ctx, q)
	if err != nil {
		return nil, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengambil data pendidikan",
		}
	}
	defer rows.Close()

	for rows.Next() {
		var m models.PendidikanModel
		err = rows.Scan(&m.ID, &m.Nama)
		if err != nil {
			return nil, err
		}
		r = append(r, m)
	}

	return r, nil
}
