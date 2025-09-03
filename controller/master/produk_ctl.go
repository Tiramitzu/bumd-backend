package controller_mst

import (
	"fmt"
	"math"
	models "microdata/kemendagri/bumd/model/master"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type ProdukController struct {
	pgxConn *pgxpool.Pool
}

func NewProdukController(pgxConn *pgxpool.Pool) *ProdukController {
	return &ProdukController{pgxConn: pgxConn}
}

func (c *ProdukController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, nama string, page int, limit int) (r []models.ProdukModel, totalCount int, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := int(claims["id_daerah"].(float64))
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM mst_produk WHERE id_bumd = $1 AND id_daerah = $2 AND deleted_by = 0
	`
	args = append(args, idBumd, idDaerah)
	if nama != "" {
		q += fmt.Sprintf(` AND nama ILIKE $%d`, len(args)+1)
		args = append(args, "%"+nama+"%")
	}
	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Produk: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, id_bumd, nama, deskripsi, gambar FROM mst_produk WHERE id_bumd = $1 AND id_daerah = $2 AND deleted_by = 0
	`
	args = append(args, idBumd, idDaerah)
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
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Produk: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var m models.ProdukModel
		err = rows.Scan(&m.ID, &m.IDBumd, &m.IDDaerah, &m.Nama, &m.Deskripsi, &m.Gambar)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Produk: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}
