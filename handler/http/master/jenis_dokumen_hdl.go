package http_mst

import (
	controller "microdata/kemendagri/bumd/controller/master"
	models "microdata/kemendagri/bumd/models/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JenisDokumenHandler struct {
	Controller *controller.JenisDokumenController
	Validate   *validator.Validate
}

func NewJenisDokumenHandler(r fiber.Router, validator *validator.Validate, controller *controller.JenisDokumenController) {
	handler := &JenisDokumenHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("jenis_dokumen")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data jenis dokumen.
//
//	@Summary		get data jenis dokumen
//	@Description	get data jenis dokumen.
//	@ID				jenis_dokumen-index
//	@Tags			Jenis dokumen
//	@Produce		json
//	@Param			nama	query		string						false	"Nama Jenis dokumen"
//	@success		200		{object}	models.JenisDokumenModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		422		{array}		utils.RequestError			"Data validation failed"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_dokumen [get]
func (h *JenisDokumenHandler) Index(c *fiber.Ctx) error {
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

// View func for get data jenis dokumen by id.
//
//	@Summary		get data jenis dokumen by id
//	@Description	get data jenis dokumen by id.
//	@ID				jenis_dokumen-view
//	@Tags			Jenis dokumen
//	@Produce		json
//	@Param			id	path		string						true	"Id untuk get data jenis dokumen"	Format(uuid)
//	@success		200	{object}	models.JenisDokumenModel	"Success"
//	@Failure		400	{object}	utils.RequestError			"Bad request"
//	@Failure		404	{object}	utils.RequestError			"Data not found"
//	@Failure		500	{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_dokumen/{id} [get]
func (h *JenisDokumenHandler) View(c *fiber.Ctx) error {
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

// Create func for create data jenis dokumen.
//
//	@Summary		create data jenis dokumen
//	@Description	create data jenis dokumen.
//	@ID				jenis_dokumen-create
//	@Tags			Jenis dokumen
//	@Accept			json
//	@Param			payload	body	models.JenisDokumenForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_dokumen [post]
func (h *JenisDokumenHandler) Create(c *fiber.Ctx) error {
	payload := new(models.JenisDokumenForm)
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

// Update func for update data jenis dokumen.
//
//	@Summary		update data jenis dokumen
//	@Description	update data jenis dokumen.
//	@ID				jenis_dokumen-update
//	@Tags			Jenis dokumen
//	@Accept			json
//	@Param			id		path	string					true	"Id untuk update data jenis dokumen"	Format(uuid)
//	@Param			payload	body	models.JenisDokumenForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_dokumen/{id} [put]
func (h *JenisDokumenHandler) Update(c *fiber.Ctx) error {
	payload := new(models.JenisDokumenForm)
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

// Delete func for delete data jenis dokumen.
//
//	@Summary		delete data jenis dokumen
//	@Description	delete data jenis dokumen.
//	@ID				jenis_dokumen-delete
//	@Tags			Jenis dokumen
//	@Accept			json
//	@Param			id	path	string	true	"Id untuk delete data jenis dokumen"	Format(uuid)
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_dokumen/{id} [delete]
func (h *JenisDokumenHandler) Delete(c *fiber.Ctx) error {
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
