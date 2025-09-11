package others

import (
	"math"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type PegawaiController struct {
	pgxConn *pgxpool.Pool
}

func NewPegawaiController(pgxConn *pgxpool.Pool) *PegawaiController {
	return &PegawaiController{pgxConn: pgxConn}
}

func (c *PegawaiController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, idBumd int, search string) (r []others.PegawaiModel, totalCount, pageCount int, err error) {
	r = make([]others.PegawaiModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_pegawai WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, id_bumd, tahun, status_pegawai, pendidikan, jumlah_pegawai
	FROM trn_pegawai
	WHERE deleted_by = 0 AND id_bumd = $1
	ORDER BY id DESC
	LIMIT $2 OFFSET $3
	`
	args = append(args, idBumd, limit, offset)

	if search != "" {
		qCount += ` AND (tahun ILIKE '%' || $2 || '%' OR status_pegawai ILIKE '%' || $2 || '%' OR pendidikan ILIKE '%' || $2 || '%' OR jumlah_pegawai ILIKE '%' || $2 || '%')`
		q += ` AND (tahun ILIKE '%' || $2 || '%' OR status_pegawai ILIKE '%' || $2 || '%' OR pendidikan ILIKE '%' || $2 || '%' OR jumlah_pegawai ILIKE '%' || $2 || '%')`
		args = append(args, search)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return
	}

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var m others.PegawaiModel
		err = rows.Scan(&m.ID, &m.IDBumd, &m.Tahun, &m.StatusPegawai, &m.Pendidikan, &m.JumlahPegawai)
		if err != nil {
			return
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *PegawaiController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r others.PegawaiModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, id_bumd, tahun, status_pegawai, pendidikan, jumlah_pegawai
	FROM trn_pegawai
	WHERE deleted_by = 0 AND id_bumd = $1 AND id = $2
	`

	err = c.pgxConn.QueryRow(fCtx, q, idBumd, id).Scan(&r.ID, &r.IDBumd, &r.Tahun, &r.StatusPegawai, &r.Pendidikan, &r.JumlahPegawai)
	if err != nil {
		return
	}

	return r, err
}

func (c *PegawaiController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *others.PegawaiForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO trn_pegawai (id_bumd, tahun, status_pegawai, pendidikan, jumlah_pegawai, created_by)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	_, err = c.pgxConn.Exec(fCtx, q, idBumd, payload.Tahun, payload.StatusPegawai, payload.Pendidikan, payload.JumlahPegawai, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat data pegawai. - " + err.Error(),
		}
	}

	return true, err
}

func (c *PegawaiController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int, payload *others.PegawaiForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_pegawai SET tahun = $1, status_pegawai = $2, pendidikan = $3, jumlah_pegawai = $4, updated_at = NOW(), updated_by = $5 WHERE id = $6 AND id_bumd = $7
	`
	_, err = c.pgxConn.Exec(fCtx, q, payload.Tahun, payload.StatusPegawai, payload.Pendidikan, payload.JumlahPegawai, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate data pegawai. - " + err.Error(),
		}
	}

	return true, err
}

func (c *PegawaiController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_pegawai SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghapus data pegawai. - " + err.Error(),
		}
	}

	return true, err
}
