package controller_mst

import (
	models "microdata/kemendagri/bumd/models/master"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type BisnisMatchingController struct {
	pgxConn *pgxpool.Pool
}

func NewBisnisMatchingController(pgxConn *pgxpool.Pool) *BisnisMatchingController {
	return &BisnisMatchingController{pgxConn: pgxConn}
}

func (c *BisnisMatchingController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token) (r []models.BisnisMatchingModel, err error) {
	r = make([]models.BisnisMatchingModel, 0)

	q := `SELECT id_bumd, nama_bumd, logo_bumd FROM trn_bumd WHERE deleted_by = 0`
	rows, err := c.pgxConn.Query(fCtx, q)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	for rows.Next() {
		var m models.BisnisMatchingModel
		err = rows.Scan(&m.IdBumd, &m.NamaBumd, &m.LogoBumd)
		if err != nil {
			return r, err
		}
		q = `SELECT id_produk, nama_produk, foto_produk, is_show FROM trn_produk WHERE id_bumd = $1 AND deleted_by = 0`
		rowsPS, err := c.pgxConn.Query(fCtx, q+" AND is_show = 1", m.IdBumd)
		if err != nil {
			return r, err
		}
		defer rowsPS.Close()
		m.ProdukShow = make([]models.ProdukModel, 0)
		for rowsPS.Next() {
			var p models.ProdukModel
			err = rowsPS.Scan(&p.IdProduk, &p.NamaProduk, &p.FotoProduk, &p.IsShow)
			if err != nil {
				return r, err
			}
			m.ProdukShow = append(m.ProdukShow, p)
		}

		rowsP, err := c.pgxConn.Query(fCtx, q+" AND is_show = 0", m.IdBumd)
		if err != nil {
			return r, err
		}
		defer rowsP.Close()
		m.Produk = make([]models.ProdukModel, 0)
		for rowsP.Next() {
			var p models.ProdukModel
			err = rowsP.Scan(&p.IdProduk, &p.NamaProduk, &p.FotoProduk, &p.IsShow)
			if err != nil {
				return r, err
			}
			m.Produk = append(m.Produk, p)
		}

		if len(m.Produk) > 0 || len(m.ProdukShow) > 0 {
			r = append(r, m)
		}
	}
	return r, err
}
