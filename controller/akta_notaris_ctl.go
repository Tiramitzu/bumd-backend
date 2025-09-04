package controller

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models"
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
) (r []models.AktaNotarisModel, totalCount, pageCount int, err error) {
	r = make([]models.AktaNotarisModel, 0)
	offset := limit * (page - 1)

	var args []interface{}
	q := `
	SELECT COALESCE(COUNT(*), 0) FROM akta_notaris WHERE deleted_by = 0 AND id_bumd = $1
	`
	args = append(args, idBumd)
	if search != "" {
		q += fmt.Sprintf(
			` AND (nomor_perda ILIKE $%d OR tanggal_perda ILIKE $%d OR keterangan ILIKE $%d OR modal_dasar ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Bentuk Usaha: %w", err)
	}

	args = make([]interface{}, 0)
	q = `
	SELECT id, nomor_perda, tanggal_perda, keterangan, file, modal_dasar, id_bumd FROM akta_notaris WHERE deleted_by = 0 AND id_bumd = $1
	`
	args = append(args, idBumd)
	if search != "" {
		q += fmt.Sprintf(
			` AND (nomor_perda ILIKE $%d OR tanggal_perda ILIKE $%d OR keterangan ILIKE $%d OR modal_dasar ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+search+"%")
	}
	q += fmt.Sprintf(`
	ORDER BY id DESC
	LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Bentuk Usaha: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var m models.AktaNotarisModel
		err = rows.Scan(&m.ID, &m.NomorPerda, &m.TanggalPerda, &m.Keterangan, &m.File, &m.ModalDasar, &m.IDBumd)
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

func (c *AktaNotarisController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r models.AktaNotarisModel, err error) {
	q := `
	SELECT id, nomor_perda, tanggal_perda, keterangan, file, modal_dasar, id_bumd FROM akta_notaris WHERE id = $1
	`

	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&r.ID, &r.NomorPerda, &r.TanggalPerda, &r.Keterangan, &r.File, &r.ModalDasar, &r.IDBumd)
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data Akta Notaris: %w", err)
	}

	return r, err
}

func (c *AktaNotarisController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.AktaNotarisForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))

	q := `
	INSERT INTO akta_notaris (nomor_perda, tanggal_perda, keterangan, file, modal_dasar, id_bumd, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.NomorPerda, payload.TanggalPerda, payload.Keterangan, payload.File, payload.ModalDasar, payload.IDBumd, idUser)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *AktaNotarisController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.AktaNotarisForm, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))

	var args []interface{}
	q := `
	UPDATE akta_notaris
	SET nomor_perda = $1, tanggal_perda = $2, keterangan = $3, file = $4, modal_dasar = $5, updated_by = $6
	WHERE id = $7
	`
	args = append(args, payload.NomorPerda, payload.TanggalPerda, payload.Keterangan, payload.File, payload.ModalDasar, idUser, id)

	_, err = c.pgxConn.Exec(fCtx, q, args...)
	if err != nil {
		return false, err
	}

	return true, err
}

func (c *AktaNotarisController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))

	q := `
	UPDATE akta_notaris
	SET deleted_by = $1, deleted_at = $2
	WHERE id = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, err
	}

	return true, err
}
