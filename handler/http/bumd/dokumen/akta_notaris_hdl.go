package dokumen

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/dokumen"
	"microdata/kemendagri/bumd/models/bumd/dokumen"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AktaNotarisHandler struct {
	Controller *ctl.AktaNotarisController
	Validate   *validator.Validate
}

func NewAktaNotarisHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *ctl.AktaNotarisController,
) {
	handler := &AktaNotarisHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("akta_notaris")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data akta notaris.
//
//	@Summary		get data akta notaris
//	@Description	get data akta notaris.
//	@ID				akta_notaris-index
//	@Tags			Akta Notaris
//	@Produce		json
//	@Param			id_bumd	path		int							true	"Id BUMD"
//	@Param			page	query		int							false	"Halaman yang ditampilkan"
//	@Param			limit	query		int							false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string						false	"Pencarian"
//	@success		200		{object}	dokumen.AktaNotarisModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		422		{array}		utils.RequestError			"Data validation failed"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/akta_notaris [get]
func (h *AktaNotarisHandler) Index(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
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
		idBumd,
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
		c.Append("x-pagination-previous-page", strconv.Itoa(page-1))
	}
	c.Append("x-pagination-current-page", strconv.Itoa(page))
	if page < pageCount {
		c.Append("x-pagination-next-page", strconv.Itoa(page+1))
	}

	return c.JSON(m)
}

// View func for get data akta notaris by id.
//
//	@Summary		get data akta notaris by id
//	@Description	get data akta notaris by id.
//	@ID				akta_notaris-view
//	@Tags			Akta Notaris
//	@Produce		json
//	@Param			id_bumd	path		int							true	"Id BUMD"
//	@Param			id		path		int							true	"Id untuk get data akta notaris"
//	@success		200		{object}	dokumen.AktaNotarisModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/akta_notaris/{id} [get]
func (h *AktaNotarisHandler) View(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create data akta notaris.
//
//	@Summary		create data akta notaris
//	@Description	create data akta notaris.
//	@ID				akta_notaris-create
//	@Tags			Akta Notaris
//	@Accept			multipart/form-data
//	@Param			id_bumd		path		int		true	"Id BUMD"
//	@Param			nomor		formData	string	true	"Nomor Akta notaris"
//	@Param			notaris		formData	string	true	"Notaris"
//	@Param			tanggal		formData	string	true	"Tanggal Akta notaris"
//	@Param			keterangan	formData	string	true	"Keterangan"
//	@Param			file		formData	file	false	"File"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/akta_notaris [post]
func (h *AktaNotarisHandler) Create(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	payload := new(dokumen.AktaNotarisForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.File != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.File, maxFileSize, []string{"application/pdf"}); err != nil {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	m, err := h.Controller.Create(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idBumd,
		payload,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}

// Update func for update data akta notaris.
//
//	@Summary		update data akta notaris
//	@Description	update data akta notaris.
//	@ID				akta_notaris-update
//	@Tags			Akta Notaris
//	@Accept			multipart/form-data
//	@Param			id_bumd		path		int		true	"Id BUMD"
//	@Param			id			path		int		true	"Id untuk update data akta notaris"
//	@Param			nomor		formData	string	true	"Nomor Akta notaris"
//	@Param			tanggal		formData	string	true	"Tanggal Akta notaris"
//	@Param			keterangan	formData	string	true	"Keterangan"
//	@Param			file		formData	file	false	"File"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/akta_notaris/{id} [put]
func (h *AktaNotarisHandler) Update(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(dokumen.AktaNotarisForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.File != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.File, maxFileSize, []string{"application/pdf"}); err != nil {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	m, err := h.Controller.Update(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idBumd,
		id,
		payload,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}

// Delete func for delete data akta notaris.
//
//	@Summary		delete data akta notaris
//	@Description	delete data akta notaris.
//	@ID				akta_notaris-delete
//	@Tags			Akta Notaris
//	@Accept			json
//	@Param			id_bumd	path	int	true	"Id BUMD"
//	@Param			id		path	int	true	"Id untuk delete data akta notaris"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/akta_notaris/{id} [delete]
func (h *AktaNotarisHandler) Delete(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	m, err := h.Controller.Delete(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idBumd,
		id,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}
