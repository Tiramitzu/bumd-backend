package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"
	models "microdata/kemendagri/bumd/models/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
//	@Param			nama	query		string					false	"Nama Bentuk Usaha"
//	@success		200		{object}	models.BentukUsahaModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha [get]
func (h *BentukUsahaHandler) Index(c *fiber.Ctx) error {
	nama := c.Query("nama")

	m, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		nama,
	)
	if err != nil {
		return err
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
//	@Param			id	path		string					true	"Id untuk get data bentuk usaha"	Format(uuid)
//	@success		200	{object}	models.BentukUsahaModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha/{id} [get]
func (h *BentukUsahaHandler) View(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.View(c.Context(), parsedId)
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
//	@Param			id		path	string					true	"Id untuk update data bentuk usaha"	Format(uuid)
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

	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.Update(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
		parsedId,
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
//	@Param			id	path	string	true	"Id untuk delete data bentuk usaha"	Format(uuid)
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_usaha/{id} [delete]
func (h *BentukUsahaHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.Delete(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedId,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}
