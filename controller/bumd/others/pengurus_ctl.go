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

type PengurusController struct {
	pgxConn *pgxpool.Pool
}

func NewPengurusController(pgxConn *pgxpool.Pool) *PengurusController {
	return &PengurusController{pgxConn: pgxConn}
}

func (c *PengurusController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, idBumd int, search string) (r []others.PengurusModel, totalCount, pageCount int, err error) {
	r = make([]others.PengurusModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_pengurus WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT id, id_bumd, jabatan_struktur, nama_pengurus, nik, alamat, deskripsi_jabatan, pendidikan_akhir, tanggal_mulai_jabatan, tanggal_akhir_jabatan, file
	FROM trn_pengurus
	WHERE deleted_by = 0 AND id_bumd = $1
	ORDER BY id DESC
	LIMIT $2 OFFSET $3
	`
	args = append(args, idBumd, limit, offset)

	if search != "" {
		qCount += ` AND (nama_pengurus ILIKE '%' || $2 || '%' OR nik ILIKE '%' || $2 || '%')`
		q += ` AND (nama_pengurus ILIKE '%' || $2 || '%' OR nik ILIKE '%' || $2 || '%')`
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
		var m others.PengurusModel
		err = rows.Scan(&m.ID, &m.IDBumd, &m.JabatanStruktur, &m.NamaPengurus, &m.NIK, &m.Alamat, &m.DeskripsiJabatan, &m.PendidikanAkhir, &m.TanggalMulaiJabatan, &m.TanggalAkhirJabatan, &m.File)
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

func (c *PengurusController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r others.PengurusModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id, id_bumd, jabatan_struktur, nama_pengurus, nik, alamat, deskripsi_jabatan, pendidikan_akhir, tanggal_mulai_jabatan, tanggal_akhir_jabatan, file
	FROM trn_pengurus
	WHERE deleted_by = 0 AND id_bumd = $1 AND id = $2
	`

	err = c.pgxConn.QueryRow(fCtx, q, idBumd, id).Scan(&r.ID, &r.IDBumd, &r.JabatanStruktur, &r.NamaPengurus, &r.NIK, &r.Alamat, &r.DeskripsiJabatan, &r.PendidikanAkhir, &r.TanggalMulaiJabatan, &r.TanggalAkhirJabatan, &r.File)
	if err != nil {
		return
	}

	return r, err
}

func (c *PengurusController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *others.PengurusForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var id int
	q := `
	INSERT INTO trn_pengurus (id_bumd, jabatan_struktur, nama_pengurus, nik, alamat, deskripsi_jabatan, pendidikan_akhir, tanggal_mulai_jabatan, tanggal_akhir_jabatan, file, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
	`
	err = tx.QueryRow(fCtx, q, idBumd, payload.JabatanStruktur, payload.NamaPengurus, payload.NIK, payload.Alamat, payload.DeskripsiJabatan, payload.PendidikanAkhir, payload.TanggalMulaiJabatan, payload.TanggalAkhirJabatan, payload.File, idUser).Scan(&id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat pengurus. - " + err.Error(),
		}
	}

	if payload.File != nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. - " + err.Error(),
			}
		}
		defer src.Close()

		objectName := "pengurus/" + fileName

		q = `
		UPDATE trn_pengurus SET file = $1 WHERE id = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(fCtx, q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *PengurusController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int, payload *others.PengurusForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	tx, err := c.pgxConn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal memulai transaksi. - " + err.Error(),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	q := `
	UPDATE trn_pengurus SET jabatan_struktur = $1, nama_pengurus = $2, nik = $3, alamat = $4, deskripsi_jabatan = $5, pendidikan_akhir = $6, tanggal_mulai_jabatan = $7, tanggal_akhir_jabatan = $8, file = $9, updated_at = NOW(), updated_by = $10 WHERE id = $11 AND id_bumd = $12
	`
	_, err = tx.Exec(fCtx, q, payload.JabatanStruktur, payload.NamaPengurus, payload.NIK, payload.Alamat, payload.DeskripsiJabatan, payload.PendidikanAkhir, payload.TanggalMulaiJabatan, payload.TanggalAkhirJabatan, payload.File, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate pengurus. - " + err.Error(),
		}
	}

	if payload.File != nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.File.Filename)

		src, err := payload.File.Open()
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. - " + err.Error(),
			}
		}
		defer src.Close()

		objectName := "pengurus/" + fileName

		q = `
		UPDATE trn_pengurus SET file = $1 WHERE id = $2 AND id_bumd = $3
		`
		_, err = tx.Exec(fCtx, q, objectName, id, idBumd)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupdate file. - " + err.Error(),
			}
		}
	}

	return true, err
}

func (c *PengurusController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE trn_pengurus SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND id_bumd = $3
	`
	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal menghapus pengurus. - " + err.Error(),
		}
	}

	return true, err
}
