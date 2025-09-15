package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"
	models "microdata/kemendagri/bumd/models/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BentukBadanHukumHandler struct {
	Controller *controller.BentukBadanHukumController
	Validate   *validator.Validate
}

func NewBentukBadanHukumHandler(r fiber.Router, validator *validator.Validate, controller *controller.BentukBadanHukumController) {
	handler := &BentukBadanHukumHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("bentuk_badan_hukum")
	rStrict.Get("/:id", handler.View)
	rStrict.Get("/", handler.Index)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data bentuk badan hukum.
//
//	@Summary		get data bentuk badan hukum
//	@Description	get data bentuk badan hukum.
//	@ID				bentuk_badan_hukum-index
//	@Tags			Bentuk Badan Hukum
//	@Produce		json
//	@Param			nama	query		string							false	"Nama Bentuk Badan Hukum"
//	@success		200		{object}	models.BentukBadanHukumModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		422		{array}		utils.RequestError				"Data validation failed"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_badan_hukum [get]
func (h *BentukBadanHukumHandler) Index(c *fiber.Ctx) error {
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

// View func for get data bentuk badan hukum by id.
//
//	@Summary		get data bentuk badan hukum by id
//	@Description	get data bentuk badan hukum by id.
//	@ID				bentuk_badan_hukum-view
//	@Tags			Bentuk Badan Hukum
//	@Produce		json
//	@Param			id	path		string							true	"Id untuk get data bentuk badan hukum"	Format(uuid)
//	@success		200	{object}	models.BentukBadanHukumModel	"Success"
//	@Failure		400	{object}	utils.RequestError				"Bad request"
//	@Failure		404	{object}	utils.RequestError				"Data not found"
//	@Failure		500	{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_badan_hukum/{id} [get]
func (h *BentukBadanHukumHandler) View(c *fiber.Ctx) error {
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

// Create func for create data bentuk badan hukum.
//
//	@Summary		create data bentuk badan hukum
//	@Description	create data bentuk badan hukum.
//	@ID				bentuk_badan_hukum-create
//	@Tags			Bentuk Badan Hukum
//	@Accept			json
//	@Param			payload	body	models.BentukBadanHukumForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_badan_hukum [post]
func (h *BentukBadanHukumHandler) Create(c *fiber.Ctx) error {
	payload := new(models.BentukBadanHukumForm)
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

// Update func for update data bentuk badan hukum.
//
//	@Summary		update data bentuk badan hukum
//	@Description	update data bentuk badan hukum.
//	@ID				bentuk_badan_hukum-update
//	@Tags			Bentuk Badan Hukum
//	@Accept			json
//	@Param			id		path	string						true	"Id untuk update data bentuk badan hukum"	Format(uuid)
//	@Param			payload	body	models.BentukBadanHukumForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_badan_hukum/{id} [put]
func (h *BentukBadanHukumHandler) Update(c *fiber.Ctx) error {
	payload := new(models.BentukBadanHukumForm)
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

// Delete func for delete data bentuk badan hukum.
//
//	@Summary		delete data bentuk badan hukum
//	@Description	delete data bentuk badan hukum.
//	@ID				bentuk_badan_hukum-delete
//	@Tags			Bentuk Badan Hukum
//	@Accept			json
//	@Param			id	path	string	true	"Id untuk delete data bentuk badan hukum"	Format(uuid)
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bentuk_badan_hukum/{id} [delete]
func (h *BentukBadanHukumHandler) Delete(c *fiber.Ctx) error {
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
