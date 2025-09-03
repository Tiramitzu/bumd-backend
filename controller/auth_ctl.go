package controller

import (
	"context"
	"errors"
	"fmt"
	models "microdata/kemendagri/bumd/model"
	"microdata/kemendagri/bumd/utils"
	"net/http"
	"time"

	"github.com/gofiber/storage/redis/v3"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	contextTimeout time.Duration
	pgxConn        *pgxpool.Pool
	// pgxConnMstData *pgxpool.Pool
	jwtManager *utils.JWTManager
	Validate   *validator.Validate
	redisCl    *redis.Storage
}

func NewAuthController(
	conn *pgxpool.Pool,
	// connMstData *pgxpool.Pool,
	timeout time.Duration,
	jm *utils.JWTManager,
	vld *validator.Validate,
	redisClient *redis.Storage,
) (controller *AuthController) {
	controller = &AuthController{
		pgxConn: conn,
		// pgxConnMstData: connMstData,
		contextTimeout: timeout,
		jwtManager:     jm,
		Validate:       vld,
		redisCl:        redisClient,
	}

	return
}

func (ac *AuthController) Login(f models.LoginForm) (token, refreshToken string, err error) {
	var q string
	var user models.User

	// Validate form input
	err = ac.Validate.Struct(f)
	if err != nil {
		return
	}

	var passUser string
	q = `SELECT id, id_daerah, id_role, password FROM "users" WHERE username = $1 and deleted_by=0`
	err = ac.pgxConn.QueryRow(context.Background(), q, f.Username).Scan(&user.IdUser, &user.IdDaerah, &user.IdRole, &passUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = utils.RequestError{
				Code:    http.StatusUnauthorized,
				Message: "invalid username or password-",
			}
			return
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passUser), []byte(f.Password)); err != nil {
		err = utils.RequestError{
			Code:    http.StatusUnauthorized,
			Message: "invalid username or password",
		}
		return
	}
	var jwtExpiredDur time.Duration
	token, refreshToken, jwtExpiredDur, err = ac.jwtManager.Generate(ac.pgxConn, user)
	if err != nil {
		return
	}

	// Check and set active session in Redis
	redisKey := fmt.Sprintf("peg:%d", user.IdUser)
	existingToken, err := ac.redisCl.Get(redisKey)
	if err != nil {
		err = utils.RequestError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check token data. - " + err.Error(),
		}
		return
	} else {
		if len(existingToken) > 0 {
			// log.Println("token: ", fmt.Sprintf("%s", existingToken))
			err = ac.redisCl.Delete(redisKey)
			if err != nil {
				err = utils.RequestError{
					Code:    http.StatusInternalServerError,
					Message: "Failed to remove token data. - " + err.Error(),
				}
				return
			}
		}

		err = ac.redisCl.Set(redisKey, []byte(token), jwtExpiredDur)
		if err != nil {
			err = utils.RequestError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to set token data. - " + err.Error(),
			}
			return
		}
	}

	return
}
