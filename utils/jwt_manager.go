package utils

import (
	"context"
	"errors"
	"fmt"
	"microdata/kemendagri/bumd/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MyCustomClaim struct {
	jwt.RegisteredClaims
	KodeDDN         string `json:"kode_ddn"`
	KodeProvinsi    string `json:"kode_provinsi"`
	SubDomainDaerah string `json:"sub_domain_daerah"`
	IdUser          int64  `json:"id_user"`
	IdDaerah        int32  `json:"id_daerah"`
	IdRole          int32  `json:"id_role"`
	IdBumd          int32  `json:"id_bumd"`
}

type JWTManager struct {
	secretKey string
	issuer    string
}

func NewJWTManager(secretKey, iss string) *JWTManager {
	return &JWTManager{secretKey, iss}
}

func (m *JWTManager) Generate(dbConn *pgxpool.Pool, user models.User) (token, refreshToken string, jwtExpDuration time.Duration, err error) {
	// ambil data durasi expired jwt dan refresh_token dari table sys config
	var jwtExpiredMinutes, refreshTokenExpiredHour int
	qStr := `SELECT "jwt_expired_minutes", "refresh_token_expired_hour" FROM "sys_config"`
	err = dbConn.QueryRow(context.Background(), qStr).Scan(&jwtExpiredMinutes, &refreshTokenExpiredHour)
	if err != nil {
		return
	}

	jwtSub := fmt.Sprintf("%d.%d.%d.%d", user.IdUser, user.IdDaerah, user.IdRole, user.IdBumd)

	jwtExpDuration = time.Duration(jwtExpiredMinutes) * time.Minute

	// Create jwt token
	jwtClaims := MyCustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   jwtSub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpDuration)),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		KodeDDN:         user.KodeDDN,
		KodeProvinsi:    user.KodeProvinsi,
		SubDomainDaerah: user.SubDomainDaerah,
		IdUser:          user.IdUser,
		IdDaerah:        user.IdDaerah,
		IdRole:          user.IdRole,
		IdBumd:          user.IdBumd,
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		return
	}

	// Create refresh token
	refreshTokenClaims := MyCustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   jwtSub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(refreshTokenExpiredHour))),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		KodeDDN:         user.KodeDDN,
		KodeProvinsi:    user.KodeProvinsi,
		SubDomainDaerah: user.SubDomainDaerah,
		IdUser:          user.IdUser,
		IdDaerah:        user.IdDaerah,
		IdRole:          user.IdRole,
		IdBumd:          user.IdBumd,
	}
	/*refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    m.issuer,
		Subject:   jwtSub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(refreshTokenExpiredHour))),
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	}*/
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		return
	}

	return
}

func (m *JWTManager) Verify(token string) (*MyCustomClaim, error) {
	var r *MyCustomClaim
	tDecoded, err := jwt.ParseWithClaims(token, &MyCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil {
		return r, err
	}

	if claims, ok := tDecoded.Claims.(*MyCustomClaim); ok && tDecoded.Valid {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return r, errors.New("token expired")
		}
		return claims, nil
	}

	return r, errors.New("invalid token")
}
