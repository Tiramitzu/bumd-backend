package controller

import (
	"context"
	models "microdata/kemendagri/bumd/model"
	"microdata/kemendagri/bumd/utils"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"

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
	SELECT users.username, users.nama, roles.nama as role
		FROM users
	LEFT JOIN roles ON roles.id = users.id_role
		WHERE users.id=$1 AND users.id_daerah=$2 AND users.id_role=$3`
	err = c.pgxConn.QueryRow(context.Background(), qStr, userId, idDaerah, idRole).
		Scan(&r.Username, &r.NamaUser, &r.NamaRole)
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
	idRole := int(claims["id_role"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	if idRole < payload.IdRole {
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
		nama,
		logo,
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

	_, err = c.pgxConn.Exec(fCtx, q, payload.Username, payload.Password, idDaerah, idRole, payload.Nama, payload.Logo, idUser)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal membuat user - " + err.Error(),
		}
	}

	return true, err
}

func (c *UserController) Update(fCtx *fasthttp.RequestCtx, user *jwt.Token, payload *models.UserForm, id int) (r bool, err error) {
	claims := user.Claims.(jwt.MapClaims)
	idUser := int(claims["id_user"].(float64))
	idRole := int(claims["id_role"].(float64))
	idDaerah := int(claims["id_daerah"].(float64))

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM users WHERE id = $1 AND id_daerah = $2
	`
	var count int
	err = c.pgxConn.QueryRow(fCtx, q, id, idDaerah).Scan(&count)
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

	q = `
	UPDATE users
	SET
		username = $1,
		password = $2,
		id_daerah = $3,
		id_role = $4,
		nama = $5,
		logo = $6,
		updated_by = $7
	`

	_, err = c.pgxConn.Exec(fCtx, q, payload.Username, payload.Password, idDaerah, idRole, payload.Nama, payload.Logo, idUser)
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
	idDaerah := int(claims["id_daerah"].(float64))

	q := `
	SELECT COALESCE(COUNT(*), 0) FROM users WHERE id = $1 AND id_daerah = $2
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
	SELECT id_role FROM users WHERE id = $1 AND id_daerah = $2
	`
	err = c.pgxConn.QueryRow(fCtx, q, id, idDaerah).Scan(&IdRole)
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
	`

	_, err = c.pgxConn.Exec(fCtx, q, idUser, time.Now(), id, idDaerah)
	if err != nil {
		return false, utils.RequestError{
			Code:    fasthttp.StatusInternalServerError,
			Message: "Gagal menghapus User - " + err.Error(),
		}
	}

	return true, err
}
