package controller

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SiteController struct {
	pgxConn        *pgxpool.Pool
	contextTimeout time.Duration
}

func NewSiteController(conn *pgxpool.Pool, tot time.Duration) *SiteController {
	return &SiteController{
		contextTimeout: tot,
	}
}

func (s *SiteController) Index() (interface{}, error) {
	var err error
	var r interface{}

	r = "index"

	return r, err
}
