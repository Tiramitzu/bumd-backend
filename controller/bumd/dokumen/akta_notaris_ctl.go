package dokumen

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/dokumen"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type AktaNotarisController struct {
	pgxConn *pgxpool.Pool
}

func NewAktaNotarisController(pgxConn *pgxpool.Pool) *AktaNotarisController {
	return &AktaNotarisController{pgxConn: pgxConn}
}

func (c *AktaNotarisController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit int,
	search string,
) (r []dokumen.AktaNotarisModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]dokumen.AktaNotarisModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM dkmn_akta_notaris WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, nomor, notaris, tanggal, keterangan, file FROM dkmn_akta_notaris WHERE deleted_by = 0 AND id_bumd = $1
	`

	args := make([]interface{}, 0)
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(
			` AND (nomor ILIKE $%d OR tanggal ILIKE $%d OR tanggal ILIKE $%d OR keterangan ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		q += fmt.Sprintf(
			` AND (nomor ILIKE $%d OR tanggal ILIKE $%d OR tanggal ILIKE $%d OR keterangan ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Akta Notaris: %w", err)
	}

	q += fmt.Sprintf(`ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m dokumen.AktaNotarisModel
		err = rows.Scan(&m.ID, &m.Nomor, &m.Notaris, &m.Tanggal, &m.Keterangan, &m.File)
		m.IDBumd = int32(idBumd)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Bentuk Usaha: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *AktaNotarisController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r dokumen.AktaNotarisModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, nomor, notaris, tanggal, keterangan, file FROM dkmn_akta_notaris WHERE id = $1 AND id_bumd = $2 AND deleted_by = 0
	`

	err = c.pgxConn.QueryRow(fCtx, q, id, idBumd).Scan(&r.ID, &r.Nomor, &r.Notaris, &r.Tanggal, &r.Keterangan, &r.File)
	r.IDBumd = int32(idBumd)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, utils.RequestError{
				Code:    fasthttp.StatusNotFound,
				Message: "Data Akta Notaris tidak ditemukan",
			}
		}
		return r, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	return r, err
}

func (c *AktaNotarisController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *dokumen.AktaNotarisForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO dkmn_akta_notaris (nomor, notaris, tanggal, keterangan, file, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Nomor, payload.Notaris, payload.Tanggal, payload.Keterangan, payload.File, idBumd, idUser)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *AktaNotarisController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int, payload *dokumen.AktaNotarisForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	q := `
	UPDATE dkmn_akta_notaris
	SET nomor = $1, notaris = $2, tanggal = $3, keterangan = $4, file = $5, updated_by = $6
	WHERE id = $7 AND id_bumd = $8
	`
	args = append(args, payload.Nomor, payload.Notaris, payload.Tanggal, payload.Keterangan, payload.File, idUser, id, idBumd)

	_, err = c.pgxConn.Exec(fCtx, q, args...)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *AktaNotarisController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE dkmn_akta_notaris
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3 AND id_bumd = $4
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id, idBumd)
	if err != nil {
		return false, err
	}

	return true, err
}
