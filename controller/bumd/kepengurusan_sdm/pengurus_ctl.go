package kepengurusan_sdm

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models/bumd/kepengurusan_sdm"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

type PengurusController struct {
	pgxConn   *pgxpool.Pool
	minioConn *utils.MinioConn
}

func NewPengurusController(pgxConn *pgxpool.Pool, minioConn *utils.MinioConn) *PengurusController {
	return &PengurusController{pgxConn: pgxConn, minioConn: minioConn}
}

func (c *PengurusController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, idBumd uuid.UUID, search string) (r []kepengurusan_sdm.PengurusModel, totalCount, pageCount int, err error) {
	r = make([]kepengurusan_sdm.PengurusModel, 0)
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

	qCount := `SELECT COALESCE(COUNT(*), 0) FROM trn_pengurus WHERE deleted_by = 0 AND id_bumd = $1`
	q := `
	SELECT
		id_pengurus,
		id_bumd,
		jabatan_struktur_pengurus,
		nama_pengurus,
		nik_pengurus,
		alamat_pengurus,
		deskripsi_jabatan_pengurus,
		pendidikan_akhir_pengurus,
		m_pendidikan.nama_pendidikan,
		tanggal_mulai_jabatan_pengurus,
		tanggal_akhir_jabatan_pengurus,
		file_pengurus,
		is_active,
		trn_pengurus.created_at,
		trn_pengurus.created_by,
		trn_pengurus.updated_at,
		trn_pengurus.updated_by
	FROM trn_pengurus
	LEFT JOIN m_pendidikan ON m_pendidikan.id_pendidikan = pendidikan_akhir_pengurus
	WHERE trn_pengurus.deleted_by = 0 AND id_bumd = $1
	`
	args = append(args, idBumd)

	if search != "" {
		qCount += fmt.Sprintf(` AND (nama_pengurus_pengurus ILIKE $%d OR nik_pengurus ILIKE $%d)`, len(args)+1, len(args)+2)
		q += fmt.Sprintf(` AND (nama_pengurus_pengurus ILIKE $%d OR nik_pengurus ILIKE $%d)`, len(args)+1, len(args)+2)
		args = append(args, search)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data pengurus. - " + err.Error(),
		}
	}

	q += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		return r, totalCount, pageCount, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data pengurus. - " + err.Error(),
		}
	}

	defer rows.Close()
	for rows.Next() {
		var m kepengurusan_sdm.PengurusModel
		err = rows.Scan(&m.Id, &m.IdBumd, &m.JabatanStruktur, &m.NamaPengurus, &m.NIK, &m.Alamat, &m.DeskripsiJabatan, &m.PendidikanAkhir, &m.NamaPendidikanAkhir, &m.TanggalMulaiJabatan, &m.TanggalAkhirJabatan, &m.File, &m.IsActive, &m.CreatedAt, &m.CreatedBy, &m.UpdatedAt, &m.UpdatedBy)
		if err != nil {
			return r, totalCount, pageCount, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengambil data pengurus. - " + err.Error(),
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

func (c *PengurusController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r kepengurusan_sdm.PengurusModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idBumdClaims, err := uuid.Parse(claims["id_bumd"].(string))
	if err != nil {
		return r, err
	}

	if idBumdClaims != uuid.Nil {
		idBumd = idBumdClaims
	}

	q := `
	SELECT id_pengurus, id_bumd, jabatan_struktur_pengurus, nama_pengurus, nik_pengurus, alamat_pengurus, deskripsi_jabatan_pengurus, pendidikan_akhir_pengurus, m_pendidikan.nama_pendidikan, tanggal_mulai_jabatan_pengurus, tanggal_akhir_jabatan_pengurus, file_pengurus, is_active, trn_pengurus.created_at, trn_pengurus.created_by, trn_pengurus.updated_at, trn_pengurus.updated_by
	FROM trn_pengurus
	LEFT JOIN m_pendidikan ON m_pendidikan.id_pendidikan = pendidikan_akhir_pengurus
	WHERE trn_pengurus.deleted_by = 0 AND id_bumd = $1 AND id_pengurus = $2
	`

	err = c.pgxConn.QueryRow(fCtx, q, idBumd, id).Scan(&r.Id, &r.IdBumd, &r.JabatanStruktur, &r.NamaPengurus, &r.NIK, &r.Alamat, &r.DeskripsiJabatan, &r.PendidikanAkhir, &r.NamaPendidikanAkhir, &r.TanggalMulaiJabatan, &r.TanggalAkhirJabatan, &r.File, &r.IsActive, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return r, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengambil data pengurus. - " + err.Error(),
		}
	}

	return r, err
}

func (c *PengurusController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd uuid.UUID, payload *kepengurusan_sdm.PengurusForm) (r bool, err error) {
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

	id, err := uuid.NewV7()
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal membuat id pengurus. - " + err.Error(),
		}
	}

	q := `
	INSERT INTO trn_pengurus (id_pengurus, id_bumd, jabatan_struktur_pengurus, nama_pengurus, nik_pengurus, alamat_pengurus, deskripsi_jabatan_pengurus, pendidikan_akhir_pengurus, tanggal_mulai_jabatan_pengurus, tanggal_akhir_jabatan_pengurus, is_active, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = tx.Exec(fCtx, q, id, idBumd, payload.JabatanStruktur, payload.NamaPengurus, payload.NIK, payload.Alamat, payload.DeskripsiJabatan, payload.PendidikanAkhir, payload.TanggalMulaiJabatan, payload.TanggalAkhirJabatan, payload.IsActive, idUser)
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

		// upload file
		objectName := "pengurus/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.File.Size,
			minio.PutObjectOptions{ContentType: payload.File.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		q = `
		UPDATE trn_pengurus SET file_pengurus = $1 WHERE id_pengurus = $2 AND id_bumd = $3
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

func (c *PengurusController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *kepengurusan_sdm.PengurusForm) (r bool, err error) {
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
	UPDATE trn_pengurus SET jabatan_struktur_pengurus = $1, nama_pengurus = $2, nik_pengurus = $3, alamat_pengurus = $4, deskripsi_jabatan_pengurus = $5, pendidikan_akhir_pengurus = $6, tanggal_mulai_jabatan_pengurus = $7, tanggal_akhir_jabatan_pengurus = $8, is_active = $9, updated_at = NOW(), updated_by = $10 WHERE id_pengurus = $11 AND id_bumd = $12
	`
	_, err = tx.Exec(fCtx, q, payload.JabatanStruktur, payload.NamaPengurus, payload.NIK, payload.Alamat, payload.DeskripsiJabatan, payload.PendidikanAkhir, payload.TanggalMulaiJabatan, payload.TanggalAkhirJabatan, payload.IsActive, idUser, id, idBumd)
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

		// upload file
		objectName := "pengurus/" + fileName
		_, err = c.minioConn.MinioClient.PutObject(
			context.Background(),
			c.minioConn.BucketName,
			objectName,
			src,
			payload.File.Size,
			minio.PutObjectOptions{ContentType: payload.File.Header.Get("Content-Type")},
		)
		if err != nil {
			return false, utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal mengupload file. - " + err.Error(),
			}
		}

		q = `
		UPDATE trn_pengurus SET file_pengurus = $1 WHERE id_pengurus = $2 AND id_bumd = $3
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

func (c *PengurusController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID) (r bool, err error) {
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
	UPDATE trn_pengurus SET deleted_by = $1, deleted_at = NOW() WHERE id_pengurus = $2 AND id_bumd = $3
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

func (c *PengurusController) UpdateStatus(fCtx *fasthttp.RequestCtx, user *jwt.Token, idBumd, id uuid.UUID, payload *kepengurusan_sdm.PengurusUpdateForm) (r bool, err error) {
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
	UPDATE trn_pengurus SET is_active = $1, updated_at = NOW(), updated_by = $2 WHERE id_pengurus = $3 AND id_bumd = $4
	`
	_, err = c.pgxConn.Exec(fCtx, q, payload.IsActive, idUser, id, idBumd)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "gagal mengupdate status pengurus. - " + err.Error(),
		}
	}

	return true, err
}
