package kepengurusan_sdm

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/kepengurusan_sdm"
	"microdata/kemendagri/bumd/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type PegawaiController struct {
	pgxConn *pgxpool.Pool
}

func NewPegawaiController(pgxConn *pgxpool.Pool) *PegawaiController {
	return &PegawaiController{pgxConn: pgxConn}
}

func (c *PegawaiController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, idBumd uuid.UUID, search string) (r []kepengurusan_sdm.PegawaiModel, totalCount, pageCount int, err error) {
	r = make([]kepengurusan_sdm.PegawaiModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, totalCount, pageCount, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	var args []interface{}
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_pegawai WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id_pegawai, id_bumd, tahun_pegawai, status_pegawai, pendidikan_pegawai, jumlah_pegawai
	FROM trn_pegawai WHERE deleted_by = 0 AND id_bumd = $1
	`
	args = append(args, idBumd)

	if search != "" {
		qCount += fmt.Sprintf(` AND (tahun_pegawai ILIKE $%d OR status_pegawai ILIKE $%d OR pendidikan_pegawai ILIKE $%d OR jumlah_pegawai ILIKE $%d)`, len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		q += fmt.Sprintf(` AND (tahun_pegawai ILIKE $%d OR status_pegawai ILIKE $%d OR pendidikan_pegawai ILIKE $%d OR jumlah_pegawai ILIKE $%d)`, len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, search)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data pegawai. - " + err.Error(),
		}
	}

	q += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data pegawai. - " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m kepengurusan_sdm.PegawaiModel
		err = rows.Scan(&m.Id, &m.IdBumd, &m.Tahun, &m.StatusPegawai, &m.Pendidikan, &m.JumlahPegawai)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data pegawai. - " + err.Error(),
			}
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *PegawaiController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r kepengurusan_sdm.PegawaiModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_pegawai, id_bumd, tahun_pegawai, status_pegawai, pendidikan_pegawai, jumlah_pegawai
	FROM trn_pegawai
	WHERE deleted_by = 0 AND id_bumd = $1 AND id_pegawai = $2
	`

	err = c.pgxConn.QueryRow(fCtx, q, idBumd, id).Scan(&r.Id, &r.IdBumd, &r.Tahun, &r.StatusPegawai, &r.Pendidikan, &r.JumlahPegawai)
	if err != nil {
		return
	}

	return r, err
}

func (c *PegawaiController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *kepengurusan_sdm.PegawaiForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id pegawai. - " + err.Error(),
		}
	}

	q := `
	INSERT INTO trn_pegawai (id_pegawai, id_bumd, tahun_pegawai, status_pegawai, pendidikan_pegawai, jumlah_pegawai, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = c.pgxConn.Exec(fCtx, q, id, idBumd, payload.Tahun, payload.StatusPegawai, payload.Pendidikan, payload.JumlahPegawai, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat data pegawai. - " + err.Error(),
		}
	}

	return true, err
}

func (c *PegawaiController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *kepengurusan_sdm.PegawaiForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_pegawai SET tahun_pegawai = $1, status_pegawai = $2, pendidikan_pegawai = $3, jumlah_pegawai = $4, updated_at = NOW(), updated_by = $5 WHERE id_pegawai = $6 AND id_bumd = $7
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

func (c *PegawaiController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id BUMD. - " + err.Error(),
		}
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_pegawai SET deleted_by = $1, deleted_at = NOW() WHERE id_pegawai = $2 AND id_bumd = $3
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
