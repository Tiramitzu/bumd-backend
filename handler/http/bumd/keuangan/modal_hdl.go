package keuangan

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/keuangan"
	"microdata/kemendagri/bumd/models/bumd/keuangan"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ModalHandler struct {
	Controller *ctl.ModalController
	Validate   *validator.Validate
}

func NewModalHandler(r fiber.Router, validator *validator.Validate, controller *ctl.ModalController) {
	handler := &ModalHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("modal")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data modal.
//
//	@Summary		get data modal
//	@Description	get data modal.
//	@ID				modal-index
//	@Tags			Modal
//	@Produce		json
//	@Param			id_bumd	path		string					true	"Id BUMD"	Format(uuid)
//	@Param			page	query		int						false	"Halaman yang ditampilkan"
//	@Param			limit	query		int						false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string					false	"Pencarian"
//	@success		200		{object}	keuangan.KeuModalModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/modal [get]
func (h *ModalHandler) Index(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	page := c.QueryInt("page", 1)
	var limit int
	limit = c.QueryInt("limit", 5)
	if limit > 5 {
		limit = 5
	}
	search := c.Query("search")

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedIdBumd,
		page,
		limit,
		search,
	)

	if err != nil {
		return err
	}

	c.Append("x-pagination-total-count", strconv.Itoa(totalCount))
	c.Append("x-pagination-page-count", strconv.Itoa(pageCount))
	c.Append("x-pagination-page-size", strconv.Itoa(limit))
	if page > 1 {
		c.Append("x-pagination-next-page", strconv.Itoa(page+1))
	}
	c.Append("x-pagination-current-page", strconv.Itoa(page))
	if page < pageCount {
		c.Append("x-pagination-prev-page", strconv.Itoa(page-1))
	}

	return c.JSON(m)
}

// View func for get data modal by id.
//
//	@Summary		get data modal by id
//	@Description	get data modal by id.
//	@ID				modal-view
//	@Tags			Modal
//	@Produce		json
//	@Param			id_bumd	path		string					true	"Id BUMD"	Format(uuid)
//	@Param			id		path		string					true	"Id Modal"	Format(uuid)
//	@success		200		{object}	keuangan.KeuModalModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/modal/{id} [get]
func (h *ModalHandler) View(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.View(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedIdBumd,
		parsedId,
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create data modal.
//
//	@Summary		create data modal
//	@Description	create data modal.
//	@ID				modal-create
//	@Tags			Modal
//	@Accept			json
//	@Param			id_bumd	path	string					true	"Id BUMD"	Format(uuid)
//	@Param			payload	body	keuangan.KeuModalForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/modal [post]
func (h *ModalHandler) Create(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	payload := new(keuangan.KeuModalForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update data modal.
//
//	@Summary		update data modal
//	@Description	update data modal.
//	@ID				modal-update
//	@Tags			Modal
//	@Accept			json
//	@Param			id_bumd	path	string					true	"Id BUMD"	Format(uuid)
//	@Param			payload	body	keuangan.KeuModalForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/modal/{id} [put]
func (h *ModalHandler) Update(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	payload := new(keuangan.KeuModalForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete data modal.
//
//	@Summary		delete data modal
//	@Description	delete data modal.
//	@ID				modal-delete
//	@Tags			Modal
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"	Format(uuid)
//	@Param			id		path		string				true	"Id Modal"	Format(uuid)
//	@success		200		{object}	boolean				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/modal/{id} [delete]
func (h *ModalHandler) Delete(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
