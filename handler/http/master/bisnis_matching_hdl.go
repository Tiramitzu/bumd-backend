package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	rStrict.Put("/:id", handler.Update)
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

// Update func for update data bisnis matching.
//
//	@Summary		update data bisnis matching
//	@Description	update data bisnis matching.
//	@ID				bisnis_matching-update
//	@Tags			Bisnis Matching
//	@Produce		json
//	@Param			id		path		string						true	"Id produk untuk update data bisnis matching"	Format(uuid)
//	@Param			status	query		int							true	"Status untuk update data bisnis matching"		Enum(0,1)
//	@Success		200		{object}	models.BisnisMatchingModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		422		{array}		utils.RequestError			"Data validation failed"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bisnis_matching/{id} [put]
func (h *BisnisMatchingHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	status := c.QueryInt("status")
	if status != 0 && status != 1 {
		return err
	}
	m, err := h.Controller.Update(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedId,
		status,
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
