package http

import (
	"microdata/kemendagri/bumd/controller"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type DashboardHandler struct {
	Controller *controller.DashboardController
	Validate   *validator.Validate
}

func NewDashboardHandler(r fiber.Router, validator *validator.Validate, controller *controller.DashboardController) {
	handler := &DashboardHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("dashboard")
	rStrict.Get("/", handler.Index)
}

// Index func for index.
//
//	@Summary		get data dashboard
//	@Description	index page.
//	@ID				dashboard-index
//	@Tags			Dashboard
//	@Produce		json
//	@Success		200	{object}	models.DashboardModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		422	{array}		utils.RequestError		"Data validation failed"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/dashboard [get]
func (h *DashboardHandler) Index(c *fiber.Ctx) error {
	r, err := h.Controller.Index(c.Context(), c.Locals("jwt").(*jwt.Token))
	if err != nil {
		return err
	}

	return c.JSON(r)
}
