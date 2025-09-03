package http

import (
	"microdata/kemendagri/bumd/controller"
	models "microdata/kemendagri/bumd/model"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type BumdHandler struct {
	Controller *controller.BumdController
	Validate   *validator.Validate
}

func NewBumdHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *controller.BumdController,
) {
	handler := &BumdHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("bumd")
	rStrict.Get("/", handler.Index)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data bumd.
//
//	@Summary		get data bumd
//	@Description	get data bumd.
//	@ID				bumd-index
//	@Tags			BUMD
//	@Produce		json
//	@Param			nama	query		string				false	"Nama BUMD"
//	@Param			page	query		int					false	"Halaman yang ditampilkan"
//	@Param			limit	query		int					false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@success		200		{object}	models.BumdModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd [get]
func (h *BumdHandler) Index(c *fiber.Ctx) error {
	nama := c.Query("nama")
	page := c.QueryInt("page", 1)
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

// Create func for create data bumd.
//
//	@Summary		create data bumd
//	@Description	create data bumd.
//	@ID				bumd-create
//	@Tags			BUMD
//	@Accept			json
//	@Param			payload	body	models.BumdForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd [post]
func (h *BumdHandler) Create(c *fiber.Ctx) error {
	payload := new(models.BumdForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Create(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}

// Update func for update data bumd.
//
//	@Summary		update data bumd
//	@Description	update data bumd.
//	@ID				bumd-update
//	@Tags			BUMD
//	@Accept			json
//	@Param			id		path	int				true	"Id untuk update data bumd"
//	@Param			payload	body	models.BumdForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [put]
func (h *BumdHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(models.BumdForm)

	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Update(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
		id,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}

// Delete func for delete data bumd.
//
//	@Summary		delete data bumd
//	@Description	delete data bumd.
//	@ID				bumd-delete
//	@Tags			BUMD
//	@Accept			json
//	@Param			id	path	int	true	"Id untuk delete data bumd"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [delete]
func (h *BumdHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	m, err := h.Controller.Delete(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		id,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}
