package keuangan

import (
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/keuangan"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type ModalController struct {
	pgxConn *pgxpool.Pool
	// pgxConnMstData *pgxpool.Pool
}

func NewModalController(pgxConn *pgxpool.Pool) *ModalController {
	return &ModalController{pgxConn: pgxConn}
}

func (c *ModalController) Index(
	fCtx *fasthttp.RequestCtx,
	user *jwt.Token,
	idBumd,
	page,
	limit int,
	search string,
) (r []keuangan.KeuModalModel, totalCount, pageCount int, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	r = make([]keuangan.KeuModalModel, 0)
	offset := limit * (page - 1)

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM keu_modal WHERE deleted_by = 0 AND id_bumd = $1`

	q := `WITH 
	data_daerah_prov_temp AS
	(
		SELECT * FROM dblink ($2,
		'
			SELECT id_daerah, nama_daerah
			FROM data.m_daerah
			WHERE
				is_deleted = 0
		')
	AS data_daerah_prov_temp
		(id_daerah INT4, nama_daerah VARCHAR)
	),
	data_daerah_kab_temp AS
	(
		SELECT * FROM dblink ($2,
		'
			SELECT id_daerah, nama_daerah
			FROM data.m_daerah
			WHERE 
				is_deleted = 0
		')
	AS data_daerah_kab_temp
		(id_daerah INT4, nama_daerah VARCHAR)
	)
	SELECT
		id,
		id_bumd,
		id_prov,
		data_daerah_prov_temp.nama_daerah,
		id_kab,
		data_daerah_kab_temp.nama_daerah,
		pemegang,
		no_ba,
		tanggal,
		jumlah,
		keterangan
	FROM keu_modal
	LEFT JOIN data_daerah_prov_temp ON keu_modal.id_prov = data_daerah_prov_temp.id_daerah
	LEFT JOIN data_daerah_kab_temp ON keu_modal.id_kab = data_daerah_kab_temp.id_daerah
	WHERE deleted_by = 0 AND id_bumd = $1
	`
	var args []interface{}
	args = append(args, idBumd)
	if search != "" {
		qCount += fmt.Sprintf(` AND pemegang ILIKE $%d OR no_ba ILIKE $%d OR tanggal ILIKE $%d OR jumlah ILIKE $%d OR keterangan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data Modal: %w", err)
	}

	args = append(args, os.Getenv("DB_SERVER_URL_MST_DATA"))

	if search != "" {
		q += fmt.Sprintf(` AND pemegang ILIKE $%d OR no_ba ILIKE $%d OR tanggal ILIKE $%d OR jumlah ILIKE $%d OR keterangan ILIKE $%d`, len(args)+1, len(args)+1, len(args)+1, len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%")
	}

	q += fmt.Sprintf(` ORDER BY id DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data Modal: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m keuangan.KeuModalModel
		err = rows.Scan(
			&m.Id,
			&m.IdBumd,
			&m.IdProv,
			&m.NamaProv,
			&m.IdKab,
			&m.NamaKab,
			&m.Pemegang,
			&m.NoBa,
			&m.Tanggal,
			&m.Jumlah,
			&m.Keterangan,
		)
		m.IdBumd = int64(idBumd)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data Modal: %w", err)
		}
		r = append(r, m)
	}

	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *ModalController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id int) (r keuangan.KeuModalModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	var args []interface{}
	args = append(args, idBumd, id, os.Getenv("DB_SERVER_URL_MST_DATA"))

	q := `WITH
	data_daerah_prov_temp AS
	(
		SELECT * FROM dblink ($3,
		'
			SELECT id_daerah, nama_daerah
			FROM data.m_daerah
			WHERE
				is_deleted = 0
		')
		AS data_daerah_prov_temp
		(id_daerah INT4, nama_daerah VARCHAR)
	),
	data_daerah_kab_temp AS
	(
		SELECT * FROM dblink ($3,
		'
			SELECT id_daerah, nama_daerah
			FROM data.m_daerah
			WHERE
				is_deleted = 0
	')
	AS data_daerah_kab_temp
	(id_daerah INT4, nama_daerah VARCHAR)
	)
	SELECT
		id,
		id_bumd,
		id_prov,
		data_daerah_prov_temp.nama_daerah,
		id_kab,
		data_daerah_kab_temp.nama_daerah,
		pemegang,
		no_ba,
		tanggal,
		jumlah,
		keterangan
	FROM keu_modal
	LEFT JOIN data_daerah_prov_temp ON keu_modal.id_prov = data_daerah_prov_temp.id_daerah
	LEFT JOIN data_daerah_kab_temp ON keu_modal.id_kab = data_daerah_kab_temp.id_daerah
	WHERE deleted_by = 0 AND id_bumd = $1 AND id = $2
	`

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&r.Id, &r.IdBumd, &r.IdProv, &r.NamaProv, &r.IdKab, &r.NamaKab, &r.Pemegang, &r.NoBa, &r.Tanggal, &r.Jumlah, &r.Keterangan)
	if err != nil {
		return r, fmt.Errorf("gagal mengambil data Modal: %w", err)
	}

	return r, err
}

func (c *ModalController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, payload *keuangan.KeuModalForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	INSERT INTO keu_modal (id_bumd, id_prov, id_kab, pemegang, no_ba, tanggal, jumlah, keterangan, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = c.pgxConn.Exec(fCtx, q, idBumd, payload.IdProv, payload.IdKab, payload.Pemegang, payload.NoBa, payload.Tanggal, payload.Jumlah, payload.Keterangan, idUser)
	if err != nil {
		return r, fmt.Errorf("gagal membuat data Modal: %w", err)
	}
	return true, err
}

func (c *ModalController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id int, payload *keuangan.KeuModalForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE keu_modal
	SET
		id_prov = $3,
		id_kab = $4,
		pemegang = $5,
		no_ba = $6,
		tanggal = $7,
		jumlah = $8,
		keterangan = $9,
		updated_by = $10,
		updated_at = NOW()
	WHERE id = $1 
		AND id_bumd = $2
		AND deleted_by = 0
	`

	_, err = c.pgxConn.Exec(fCtx, q, id, idBumd, payload.IdProv, payload.IdKab, payload.Pemegang, payload.NoBa, payload.Tanggal, payload.Jumlah, payload.Keterangan, idUser)
	if err != nil {
		return r, fmt.Errorf("gagal mengubah data Modal: %w", err)
	}
	return true, err
}

func (c *ModalController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd int, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idBumdClaims := int(claims["id_bumd"].(float64))

	if idBumdClaims > 0 {
		idBumd = idBumdClaims
	}

	q := `
	UPDATE keu_modal
	SET
		deleted_by = $1,
		deleted_at = NOW()
	WHERE id = $2
		AND id_bumd = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, id, idBumd)
	if err != nil {
		return r, fmt.Errorf("gagal menghapus data Modal: %w", err)
	}
	return true, err
}
