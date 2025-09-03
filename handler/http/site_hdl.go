package http

import (
	"microdata/kemendagri/bumd/controller"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SiteHandler struct {
	Controller *controller.SiteController
	Validate   *validator.Validate
}

func NewSiteHandler(app *fiber.App, controller *controller.SiteController, vld *validator.Validate) {
	handler := &SiteHandler{
		Controller: controller,
		Validate:   vld,
	}

	// public route
	rPub := app.Group("/site")
	rPub.Get("/index", handler.Index)
}

// Index func for index.
//
//	@Summary		index
//	@Description	index page.
//	@ID				index
//	@Tags			Site
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	interface{}			"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		403	{object}	utils.LoginError	"Login forbidden"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Router			/site/index [get]
func (h *SiteHandler) Index(c *fiber.Ctx) error {
	r, err := h.Controller.Index()
	if err != nil {
		return err
	}

	return c.JSON(r)
}
