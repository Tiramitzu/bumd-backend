package controller_mst

import (
	"context"
	models "microdata/kemendagri/bumd/models/master"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JenisLaporanController struct {
	pgxConn *pgxpool.Pool
}

func NewJenisLaporanController(pgxConn *pgxpool.Pool) *JenisLaporanController {
	return &JenisLaporanController{pgxConn: pgxConn}
}

func (c *JenisLaporanController) Index(ctx context.Context, jwt *jwt.Token, bentukUsaha uuid.UUID, parentId int) ([]models.JenisLaporanModel, error) {
	var r []models.JenisLaporanModel

	var args []interface{}
	q := `SELECT
		id_jenis_laporan,
		bentuk_usaha,
		kode_jenis_laporan,
		uraian_jenis_laporan,
		keterangan_jenis_laporan,
		level_jenis_laporan,
		parent_id_jenis_laporan
		FROM m_jenis_laporan WHERE parent_id_jenis_laporan = $1`
	args = append(args, parentId)

	if bentukUsaha != uuid.Nil {
		q += ` AND bentuk_usaha = $2`
		args = append(args, bentukUsaha)
	}

	rows, err := c.pgxConn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.JenisLaporanModel
		err = rows.Scan(&m.Id, &m.BentukUsaha, &m.Kode, &m.Uraian, &m.Keterangan, &m.Level, &m.ParentId)
		if err != nil {
			return nil, err
		}
		r = append(r, m)
	}

	return r, nil
}
