package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type BisnisMatchingHandler struct {
	Controller *controller.BisnisMatchingController
	Validate   *validator.Validate
}

func NewBisnisMatchingHandler(r fiber.Router, validator *validator.Validate, controller *controller.BisnisMatchingController) {
	handler := &BisnisMatchingHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("bisnis_matching")
	rStrict.Get("/", handler.Index)
}

// Index func for get data bisnis matching.
//
//	@Summary		get data bisnis matching
//	@Description	get data bisnis matching.
//	@ID				bisnis_matching-index
//	@Tags			Bisnis Matching
//	@Produce		json
//	@Success		200	{object}	[]models.BisnisMatchingModel	"Success"
//	@Failure		400	{object}	utils.RequestError				"Bad request"
//	@Failure		404	{object}	utils.RequestError				"Data not found"
//	@Failure		422	{array}		utils.RequestError				"Data validation failed"
//	@Failure		500	{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bisnis_matching [get]
func (h *BisnisMatchingHandler) Index(c *fiber.Ctx) error {
	m, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
