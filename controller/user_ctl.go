package controller

import (
	"context"
	"fmt"
	"math"
	"microdata/kemendagri/bumd/models"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	contextTimeout time.Duration
	pgxConn        *pgxpool.Pool
	redisCl        *redis.Storage
}

func NewUserController(conn *pgxpool.Pool, tot time.Duration, redisClient *redis.Storage) (controller *UserController) {
	controller = &UserController{
		pgxConn:        conn,
		contextTimeout: tot,
		redisCl:        redisClient,
	}

	return
}

func (c *UserController) Index(fCtx *fasthttp.RequestCtx, user *jwt.Token, page, limit int, IdRoleFilter int) (r []models.UserModel, totalCount, pageCount int, err error) {
	r = make([]models.UserModel, 0)
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))
	offset := limit * (page - 1)

	if idRole > 2 {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki akses untuk melihat data User",
		}
		return r, totalCount, pageCount, err
	}

	var args []interface{}
	qCount := `
	SELECT COALESCE(COUNT(*), 0) FROM users WHERE deleted_by = 0
	`

	q := `
	SELECT
		users.id,
		users.id_daerah,
		users.id_role,
		users.id_bumd,
		users.username,
		users.nama,
		users.logo,
		roles.nama as role
	FROM users
	LEFT JOIN roles ON roles.id = users.id_role
	WHERE users.deleted_by = 0
	`
	if IdRoleFilter > idRole {
		qCount += fmt.Sprintf(` AND users.id_role = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND users.id_role = $%d`, len(args)+1)
		args = append(args, IdRoleFilter)
	}
	if IdRoleFilter == idRole && idRole != 1 {
		qCount += fmt.Sprintf(` AND users.id = $%d`, len(args)+1)
		q += fmt.Sprintf(` AND users.id = $%d`, len(args)+1)
		args = append(args, idUser)
	}

	err = c.pgxConn.QueryRow(fCtx, qCount, args...).Scan(&totalCount)
	if err != nil {
		return r, totalCount, pageCount, fmt.Errorf("gagal menghitung total data User: %w", err)
	}

	q += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := c.pgxConn.Query(fCtx, q, args...)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, totalCount, pageCount, fmt.Errorf("data User tidak ditemukan")
		}
		return r, totalCount, pageCount, fmt.Errorf("gagal mengambil data User: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m models.UserModel
		err = rows.Scan(&m.IdUser, &m.IdDaerah, &m.IdRole, &m.IdBumd, &m.Username, &m.Nama, &m.Logo, &m.Role)
		if err != nil {
			return r, totalCount, pageCount, fmt.Errorf("gagal memindahkan data User: %w", err)
		}
		r = append(r, m)
	}

	// page info
	pageCount = 1
	if totalCount > 0 && totalCount > limit {
		pageCount = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	return r, totalCount, pageCount, err
}

func (c *UserController) View(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r models.UserModel, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))
	if idRole > 2 {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki akses untuk melihat data User",
		}
		return r, err
	}

	var args []interface{}
	q := `
	SELECT
		users.id,
		users.id_daerah,
		users.id_role,
		users.id_bumd,
		users.username,
		users.nama,
		users.logo,
		roles.nama as role
	FROM users
	LEFT JOIN roles ON roles.id = users.id_role
	WHERE users.id = $1 AND users.id_daerah = $2 AND users.deleted_by = 0`
	args = append(args, id, idDaerah)

	if idRole > 1 && id != idUser {
		q += fmt.Sprintf(` AND users.id_role < $%d`, len(args)+1)
		args = append(args, idRole)
	}

	err = c.pgxConn.QueryRow(fCtx, q, args...).Scan(&r.IdUser, &r.IdDaerah, &r.IdRole, &r.IdBumd, &r.Username, &r.Nama, &r.Logo, &r.Role)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return r, fmt.Errorf("data User tidak ditemukan")
		}
		return r, fmt.Errorf("gagal mengambil data User: %w", err)
	}

	return r, err
}

func (c *UserController) Logout(idUser string) error {
	var err error
	err = c.redisCl.Delete("usr:" + idUser)

	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Failed to remove token data. - " + err.Error(),
		}
		return err
	}

	return err
}

/*
Profile melihat detail profile user dengan id wilayah yang sama dengan si pembuat user baru
route name auth-service/user/profile
*/
func (c *UserController) Profile(user *jwt.Token) (r models.UserDetail, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idDaerah := uint64(claims["id_daerah"].(float64))
	idRole := int(claims["id_role"].(float64))
	userId := uint64(claims["id_user"].(float64))

	qStr := `
	SELECT
		users.username,
		users.nama,
		roles.nama as role,
		COALESCE(bumd.nama, '-') as nama_bumd
	FROM users
	LEFT JOIN roles ON roles.id = users.id_role
	LEFT JOIN bumd ON bumd.id = users.id_bumd
		WHERE users.id=$1 AND users.id_daerah=$2 AND users.id_role=$3`
	err = c.pgxConn.QueryRow(context.Background(), qStr, userId, idDaerah, idRole).
		Scan(&r.Username, &r.NamaUser, &r.NamaRole, &r.NamaBumd)
	if err != nil {
		err = utils.RequestError{
			Code:    fasthttp.StatusNotFound,
			Message: err.Error(),
		}
		return
	}

	return
}

func (c *UserController) Create(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.UserForm) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int32(claims["id_role"].(float64))
	idDaerah := int32(claims["id_daerah"].(float64))

	if idDaerah < 1 {
		idDaerah = payload.IdDaerah
	}

	if payload.IdRole < 3 {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "User dengan role ini tidak dapat dibuat melalui API ini",
		}
		return
	}

	if idRole <= payload.IdRole {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki akses untuk membuat user dengan role ini",
		}
		return
	}

	q := `
	INSERT INTO users (
		username,
		password,
		id_daerah,
		id_role,
		id_bumd,
		nama,
		created_by
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
	`

	pHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengenkripsi password - " + err.Error(),
		}
	}

	_, err = c.pgxConn.Exec(fCtx, q, payload.Username, pHash, idDaerah, idRole, payload.IdBumd, payload.Nama, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat user - " + err.Error(),
		}
	}

	if payload.Logo != nil {
		// generate nama file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.Logo.Filename)

		src, err := payload.Logo.Open()
		if err != nil {
			err = utils.RequestError{
				Code:    fasthttp.StatusInternalServerError,
				Message: "gagal membuka file. " + err.Error(),
			}
			return false, err
		}
		defer src.Close()

		// upload file
		// objectName := "user/logo/" + fileName
		// _, err = c.minioConn.MinioClient.PutObject(
		// 	context.Background(),
		// 	c.minioConn.BucketName,
		// 	objectName,
		// 	src,
		// 	payload.File.Size,
		// 	minio.PutObjectOptions{ContentType: "application/pdf"},
		// )

		// if err != nil {
		// 	err = utils.RequestError{
		// 		Code:    http.StatusInternalServerError,
		// 		Message: "gagal upload file. - " + err.Error(),
		// 	}
		// 	return err
		// }
		// update file
		q = `UPDATE users SET logo=$1 WHERE id=$2`
		_, err = c.pgxConn.Exec(fCtx, q, fileName, idUser)
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

func (c *UserController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.UserForm, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int32(claims["id_role"].(float64))

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM users WHERE id = $1 AND deleted_by = 0
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghitung data User - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "User tidak ditemukan.",
		}
	}

	if idRole < payload.IdRole {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki akses untuk mengubah user dengan role ini",
		}
		return
	}

	pHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengenkripsi password - " + err.Error(),
		}
	}

	q = `
	UPDATE users
	SET
		username = $1,
		password = $2,
		id_role = $3,
		id_bumd = $4,
		nama = $5,
		logo = $6,
		updated_by = $7
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Username, pHash, idRole, payload.IdBumd, payload.Nama, payload.Logo, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal mengubah User - " + err.Error(),
		}
	}

	return true, err
}

func (c *UserController) Delete(fCtx *fasthttp.RequestCtx, user *jwt.Token, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM users WHERE id = $1
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&count)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus data User - " + err.Error(),
		}
	}

	if count < 1 {
		return false, utils.RequestError{
			Code:    fasthttp.StatusBadRequest,
			Message: "User tidak ditemukan.",
		}
	}

	var IdRole int
	q = `
	SELECT id_role FROM users WHERE id = $1
	`
	err = c.pgxConn.QueryRow(fCtx, q, id).Scan(&IdRole)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus User - " + err.Error(),
		}
	}

	if idRole < IdRole {
		err = utils.RequestError{
			Code:    fasthttp.StatusForbidden,
			Message: "Anda tidak memiliki akses untuk menghapus User dengan role ini",
		}
		return
	}

	q = `
	UPDATE users
	SET
		deleted_by = $1,
		deleted_at = $2
	WHERE id = $3
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus User - " + err.Error(),
		}
	}

	return true, err
}
