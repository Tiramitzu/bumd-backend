package controller_mst

import (
	models "microdata/kemendagri/bumd/models/master"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	claims := user.Claims.(jwt.MapClaims)
	idRole := int32(claims["id_role"].(float64))
	if idRole != 1 {
		return r, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk melihat data Bisnis Matching.",
		}
	}
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
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "Gagal mengambil data BUMD Bisnis Matching - " + err.Error(),
			}
		}
		q = `SELECT id_produk, nama_produk, foto_produk, is_show FROM trn_produk WHERE id_bumd = $1 AND deleted_by = 0`
		rowsPS, err := c.pgxConn.Query(fCtx, q+" AND is_show = 1", m.IdBumd)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "Gagal mengambil data Produk Show Bisnis Matching - " + err.Error(),
			}
		}
		defer rowsPS.Close()
		m.ProdukShow = make([]models.ProdukModel, 0)
		for rowsPS.Next() {
			var p models.ProdukModel
			err = rowsPS.Scan(&p.IdProduk, &p.NamaProduk, &p.FotoProduk, &p.IsShow)
			if err != nil {
				return r, utils.RequestError{
					Code:    fasthttp.StatusInternalServerError,
					Message: "Gagal memindahkan data Produk Bisnis Matching - " + err.Error(),
				}
			}
			m.ProdukShow = append(m.ProdukShow, p)
		}

		rowsP, err := c.pgxConn.Query(fCtx, q+" AND is_show = 0", m.IdBumd)
		if err != nil {
			return r, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "Gagal mengambil data Produk Bisnis Matching - " + err.Error(),
			}
		}
		defer rowsP.Close()
		m.Produk = make([]models.ProdukModel, 0)
		for rowsP.Next() {
			var p models.ProdukModel
			err = rowsP.Scan(&p.IdProduk, &p.NamaProduk, &p.FotoProduk, &p.IsShow)
			if err != nil {
				return r, utils.RequestError{
					Code:    fasthttp.StatusInternalServerError,
					Message: "Gagal memindahkan data Bisnis Matching - " + err.Error(),
				}
			}
			m.Produk = append(m.Produk, p)
		}

		if len(m.Produk) > 0 || len(m.ProdukShow) > 0 {
			r = append(r, m)
		}
	}
	return r, err
}

func (c *BisnisMatchingController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, id uuid.UUID, status int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int32(claims["id_role"].(float64))
	if idRole != 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki izin untuk memperbarui data Bisnis Matching.",
		}
	}

	q := `UPDATE trn_produk SET is_show = $1, updated_by = $2, updated_at = NOW() WHERE id_produk = $3`
	_, err = c.pgxConn.Exec(fCtx, q, status, idUser, id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal memperbarui data Bisnis Matching - " + err.Error(),
		}
	}
	return true, err
}
