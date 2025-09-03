package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"
	models "microdata/kemendagri/bumd/model/master"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type BentukUsahaHandler struct {
	Controller *controller.BentukUsahaController
	Validate   *validator.Validate
}

func NewBentukUsahaHandler(r fiber.Router, validator *validator.Validate, controller *controller.BentukUsahaController) {
	handler := &BentukUsahaHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("bentuk_usaha")
	rStrict.Get("/:id", handler.View)
	rStrict.Get("/", handler.Index)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data bentuk usaha.
//
//	@Summary		get data bentuk usaha
//	@Description	get data bentuk usaha.
//	@ID				bentuk_usaha-index
//	@Tags			Bentuk Usaha
//	@Produce		json
//	@Param			page	query		int						false	"Halaman yang ditampilkan"
//	@Param			limit	query		int						false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			nama	query		string					false	"Nama Bentuk Usaha"
//	@success		200		{object}	models.BentukUsahaModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha [get]
func (h *BentukUsahaHandler) Index(c *fiber.Ctx) error {
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

// View func for get data bentuk usaha by id.
//
//	@Summary		get data bentuk usaha by id
//	@Description	get data bentuk usaha by id.
//	@ID				bentuk_usaha-view
//	@Tags			Bentuk Usaha
//	@Produce		json
//	@Param			id	path		int						true	"Id untuk get data bentuk usaha"
//	@success		200	{object}	models.BentukUsahaModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha/{id} [get]
func (h *BentukUsahaHandler) View(c *fiber.Ctx) error {
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

// Create func for create data bentuk usaha.
//
//	@Summary		create data bentuk usaha
//	@Description	create data bentuk usaha.
//	@ID				bentuk_usaha-create
//	@Tags			Bentuk Usaha
//	@Accept			json
//	@Param			payload	body	models.BentukUsahaForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha [post]
func (h *BentukUsahaHandler) Create(c *fiber.Ctx) error {
	payload := new(models.BentukUsahaForm)
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

// Update func for update data bentuk usaha.
//
//	@Summary		update data bentuk usaha
//	@Description	update data bentuk usaha.
//	@ID				bentuk_usaha-update
//	@Tags			Bentuk Usaha
//	@Accept			json
//	@Param			id		path	int						true	"Id untuk update data bentuk usaha"
//	@Param			payload	body	models.BentukUsahaForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha/{id} [put]
func (h *BentukUsahaHandler) Update(c *fiber.Ctx) error {
	payload := new(models.BentukUsahaForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
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

// Delete func for delete data bentuk usaha.
//
//	@Summary		delete data bentuk usaha
//	@Description	delete data bentuk usaha.
//	@ID				bentuk_usaha-delete
//	@Tags			Bentuk Usaha
//	@Accept			json
//	@Param			id	path	int	true	"Id untuk delete data bentuk usaha"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha/{id} [delete]
func (h *BentukUsahaHandler) Delete(c *fiber.Ctx) error {
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
