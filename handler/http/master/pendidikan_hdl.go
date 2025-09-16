package http_mst

import (
	controller_mst "microdata/kemendagri/bumd/controller/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type PendidikanHandler struct {
	Controller *controller_mst.PendidikanController
	Validate   *validator.Validate
}

func NewPendidikanHandler(r fiber.Router, validator *validator.Validate, controller *controller_mst.PendidikanController) {
	handler := &PendidikanHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("pendidikan")
	rStrict.Get("/", handler.Index)
}

// Index func for get all pendidikan.
//
//	@Summary		get all pendidikan
//	@Description	get all pendidikan.
//	@ID				pendidikan-index
//	@Tags			Pendidikan
//	@Produce		json
//	@Success		200	{array}		models.PendidikanModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		422	{array}		utils.RequestError		"Data validation failed"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/pendidikan [get]
func (h *PendidikanHandler) Index(c *fiber.Ctx) error {
	r, err := h.Controller.Index(c.Context(), c.Locals("jwt").(*jwt.Token))
	if err != nil {
		return err
	}

	return c.JSON(r)
}
