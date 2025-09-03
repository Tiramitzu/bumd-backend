package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type RolesHandler struct {
	Controller *controller.RolesController
	Validate   *validator.Validate
}

func NewRolesHandler(r fiber.Router, validator *validator.Validate, controller *controller.RolesController) {
	handler := &RolesHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("roles")
	rStrict.Get("/:id", handler.View)
	rStrict.Get("/", handler.Index)
}

// Index func for get data roles.
//
//	@Summary		get data roles
//	@Description	get data roles.
//	@ID				roles-index
//	@Tags			Roles
//	@Produce		json
//	@Param			page	query		int					false	"Halaman yang ditampilkan"
//	@Param			limit	query		int					false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			nama	query		string				false	"Nama Roles"
//	@success		200		{object}	models.RolesModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/roles [get]
func (h *RolesHandler) Index(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	nama := c.Query("nama")
	var limit int
	limit = c.QueryInt("limit", 5)

	if limit > 5 {
		limit = 5
	}

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		page,
		limit,
		nama,
	)
	if err != nil {
		return err
	}

	c.Append("x-pagination-total-count", strconv.Itoa(totalCount))
	c.Append("x-pagination-page-count", strconv.Itoa(pageCount))
	c.Append("x-pagination-page-size", strconv.Itoa(limit))
	if page > 1 {
		c.Append("x-pagination-previous-page", strconv.Itoa(page-1))
	}
	c.Append("x-pagination-current-page", strconv.Itoa(page))
	if page < pageCount {
		c.Append("x-pagination-next-page", strconv.Itoa(page+1))
	}
	return c.JSON(m)
}

// View func for get data roles by id.
//
//	@Summary		get data roles by id
//	@Description	get data roles by id.
//	@ID				roles-view
//	@Tags			Bentuk Usaha
//	@Produce		json
//	@Param			id	path		int					true	"Id untuk get data roles"
//	@success		200	{object}	models.RolesModel	"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/roles/{id} [get]
func (h *RolesHandler) View(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	m, err := h.Controller.View(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
