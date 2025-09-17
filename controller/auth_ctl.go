package controller

import (
	"context"
	"errors"
	"fmt"
	"microdata/kemendagri/bumd/models"
	"microdata/kemendagri/bumd/utils"
	"net/http"
	"strings"
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
	pgxConnMstData *pgxpool.Pool
	jwtManager     *utils.JWTManager
	Validate       *validator.Validate
	redisCl        *redis.Storage
}

func NewAuthController(
	conn *pgxpool.Pool,
	connMstData *pgxpool.Pool,
	timeout time.Duration,
	jm *utils.JWTManager,
	vld *validator.Validate,
	redisClient *redis.Storage,
) (controller *AuthController) {
	controller = &AuthController{
		pgxConn:        conn,
		pgxConnMstData: connMstData,
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
		return "", "", utils.RequestError{
			Code:    http.StatusUnprocessableEntity,
			Message: "gagal validasi data - " + err.Error(),
		}
	}

	var passUser string
	q = `
	SELECT 
		id,
		id_daerah,
		id_role,
		id_bumd,
		password
	FROM "users"
	WHERE username = $1
		AND deleted_by=0
	`
	err = ac.pgxConn.QueryRow(context.Background(), q, f.Username).Scan(&user.IdUser, &user.IdDaerah, &user.IdRole, &user.IdBumd, &passUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", utils.RequestError{
				Code:    http.StatusUnauthorized,
				Message: "invalid username or password-",
			}
		}
		return "", "", utils.RequestError{
			Code:    http.StatusInternalServerError,
			Message: "gagal mengambil data user - " + err.Error(),
		}
	}

	if user.IdDaerah > 0 {
		q = `
		SELECT sub_domain, kode_ddn
		FROM data.m_daerah
		WHERE id_daerah = $1
	`
		err = ac.pgxConnMstData.QueryRow(context.Background(), q, user.IdDaerah).Scan(&user.SubDomainDaerah, &user.KodeDDN)
		if err != nil {
			return "", "", utils.RequestError{
				Code:    http.StatusInternalServerError,
				Message: "Data Pemerintah Daerah Tidak Tersedia. - " + err.Error(),
			}
		}

		user.KodeProvinsi = strings.Split(user.KodeDDN, ".")[0]
		if len(user.KodeDDN) == 2 {
			user.KodeDDN += ".00"
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passUser), []byte(f.Password)); err != nil {
		return "", "", utils.RequestError{
			Code:    http.StatusUnauthorized,
			Message: "invalid username or password",
		}
	}
	var jwtExpiredDur time.Duration
	token, refreshToken, jwtExpiredDur, err = ac.jwtManager.Generate(ac.pgxConn, user)
	if err != nil {
		return "", "", utils.RequestError{
			Code:    http.StatusInternalServerError,
			Message: "gagal menghasilkan token - " + err.Error(),
		}
	}

	// Check and set active session in Redis
	redisKey := fmt.Sprintf("usr:%d", user.IdUser)
	existingToken, err := ac.redisCl.Get(redisKey)
	if err != nil {
		return "", "", utils.RequestError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check token data. - " + err.Error(),
		}
	} else {
		if len(existingToken) > 0 {
			// log.Println("token: ", fmt.Sprintf("%s", existingToken))
			err = ac.redisCl.Delete(redisKey)
			if err != nil {
				return "", "", utils.RequestError{
					Code:    http.StatusInternalServerError,
					Message: "Failed to remove token data. - " + err.Error(),
				}
			}
		}

		err = ac.redisCl.Set(redisKey, []byte(token), jwtExpiredDur)
		if err != nil {
			return "", "", utils.RequestError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to set token data. - " + err.Error(),
			}
		}
	}

	return token, refreshToken, err
}
