package others

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type ProdukController struct {
	pgxConn *pgxpool.Pool
}

func NewProdukController(pgxConn *pgxpool.Pool) *ProdukController {
	return &ProdukController{pgxConn: pgxConn}
}

func (c *ProdukController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit int,
	search string,
) (r []others.ProdukModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]others.ProdukModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_produk WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_produk, id_bumd, nama_produk, deskripsi, foto_produk
	FROM trn_produk
	WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND nama_produk ILIKE $%d OR deskripsi ILIKE $%d`, len(args)+1, len(args)+1)
		q += fmt.Sprintf(` AND nama_produk ILIKE $%d OR deskripsi ILIKE $%d`, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data PRODUK: %w", err)
	}

	q += fmt.Sprintf(` ORDER BY id_produk DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data PRODUK: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m others.ProdukModel
		err = rows.Scan(&m.ID, &m.IDBumd, &m.NamaProduk, &m.Deskripsi, &m.FotoProduk)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data PRODUK: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *ProdukController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string) (r others.ProdukModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_produk, id_bumd, nama_produk, deskripsi, foto_produk
	FROM trn_produk
	WHERE id_produk = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.ID, &r.IDBumd, &r.NamaProduk, &r.Deskripsi, &r.FotoProduk)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data PRODUK tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data PRODUK: %w", err)
	}

	return r, err
}

func (c *ProdukController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *others.ProdukForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO trn_produk (nama_produk, deskripsi, id_bumd, created_by) VALUES ($1, $2, $3, $4) RETURNING id_produk
	`

	var id int
	err = tx.QueryRow(context.Background(), q, payload.NamaProduk, payload.Deskripsi, idBumd, idUser).Scan(&id)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memasukkan data PRODUK. - " + err.Error(),
		}
		return false, err
	}

	if payload.FotoProduk != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.FotoProduk.Filename)

		src, err := payload.FotoProduk.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "trn_produk/" + fileName

		// update file
		q = `UPDATE trn_produk SET foto_produk=$1 WHERE id_produk=$2 AND id_bumd=$3`
		_, err = tx.Exec(context.Background(), q, objectName, id, idBumd)
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
			return false, err
		}
	}

	return true, err
}

func (c *ProdukController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string, payload *others.ProdukForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	q := `
	UPDATE trn_produk
	SET nama_produk = $1, deskripsi = $2, updated_by = $3, updated_at = NOW()
	WHERE id_produk = $4 AND id_bumd = $5
	`
	args = append(args, payload.NamaProduk, payload.Deskripsi, idUser, id, idBumd)
	_, err = tx.Exec(context.Background(), q, args...)
	if err != nil {
		return false, err
	}

	if payload.FotoProduk != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.FotoProduk.Filename)

		src, err := payload.FotoProduk.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		objectName := "trn_produk/" + fileName

		// update file
		q = `UPDATE trn_produk SET foto_produk=$1 WHERE id_produk=$2`
		_, err = tx.Exec(context.Background(), q, objectName, id)
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
			return false, err
		}
	}

	return true, err
}

func (c *ProdukController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id string) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}
	q := `
	UPDATE trn_produk
	SET deleted_by = $1, deleted_at = NOW()
	WHERE id_produk = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(context.Background(), q, idUser, id, idBumd)
	if err != nil {
		return false, err
	}
	return true, err
}
