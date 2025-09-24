package controller_mst

import (
	"context"
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type JenisLaporanController struct {
	pgxConn *pgxpool.Pool
}

func NewJenisLaporanController(pgxConn *pgxpool.Pool) *JenisLaporanController {
	return &JenisLaporanController{pgxConn: pgxConn}
}

func (c *JenisLaporanController) Index(ctx context.Context, jwt *jwt.Token, bentukUsaha uuid.UUID, parentId int) ([]models.JenisLaporanModel, error) {
	r := make([]models.JenisLaporanModel, 0)
	var totalCount int
	var args []interface{}
	qCount := `SELECT COUNT(*) FROM m_jenis_laporan WHERE parent_id_jenis_laporan = $1`
	args = append(args, parentId)
	if bentukUsaha != uuid.Nil {
		qCount += ` AND bentuk_usaha = $2`
		args = append(args, bentukUsaha)
	}

	err := c.pgxConn.QueryRow(ctx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghitung total data Jenis Laporan: " + err.Error(),
		}
	}

	args = make([]interface{}, 0)
	q := `SELECT
		id_jenis_laporan,
		bentuk_usaha,
		kode_jenis_laporan,
		uraian_jenis_laporan,
		keterangan_jenis_laporan,
		level_jenis_laporan,
		parent_id_jenis_laporan
		FROM m_jenis_laporan WHERE parent_id_jenis_laporan = $1 ORDER BY created_at ASC`
	args = append(args, parentId)

	if bentukUsaha != uuid.Nil && totalCount > 0 {
		q += ` AND bentuk_usaha = $2`
		args = append(args, bentukUsaha)
	} else {
		q += ` AND bentuk_usaha = '01994c01-c285-73b3-ada9-10d5180a4a2f'`
	}

	rows, err := c.pgxConn.Query(ctx, q, args...)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data Jenis Laporan: " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		var m models.JenisLaporanModel
		err = rows.Scan(&m.Id, &m.BentukUsaha, &m.Kode, &m.Uraian, &m.Keterangan, &m.Level, &m.ParentId)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal memindahkan data Jenis Laporan: " + err.Error(),
			}
		}
		r = append(r, m)
	}

	return r, err
}
