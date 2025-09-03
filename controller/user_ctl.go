package controller

import (
	"context"
	models "microdata/kemendagri/bumd/model"
	"microdata/kemendagri/bumd/utils"
	"net/http"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-jwt/jwt/v4"

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
			Code:    http.StatusInternalServerError,
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
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		return
	}

	return
}
