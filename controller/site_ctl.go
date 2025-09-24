package controller

import (
	"time"
)

type SiteController struct {
	contextTimeout time.Duration
}

func NewSiteController(tot time.Duration) *SiteController {
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
