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

type PerdaPendirianHandler struct {
	Controller *ctl.PerdaPendirianController
	Validate   *validator.Validate
}

func NewPerdaPendirianHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *ctl.PerdaPendirianController,
) {
	handler := &PerdaPendirianHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("perda_pendirian")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data perda pendirian.
//
//	@Summary		get data perda pendirian
//	@Description	get data perda pendirian.
//	@ID				perda_pendirian-index
//	@Tags			Perda Pendirian
//	@Produce		json
//	@Param			id_bumd			path		int							true	"Id BUMD"
//	@Param			page			query		int							false	"Halaman yang ditampilkan"
//	@Param			limit			query		int							false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search			query		string						false	"Pencarian"
//	@Param			modal_dasar_min	query		float64						false	"Modal Dasar Min"
//	@Param			modal_dasar_max	query		float64						false	"Modal Dasar Max"
//	@success		200				{object}	dokumen.PerdaPendirianModel	"Success"
//	@Failure		400				{object}	utils.RequestError			"Bad request"
//	@Failure		404				{object}	utils.RequestError			"Data not found"
//	@Failure		422				{array}		utils.RequestError			"Data validation failed"
//	@Failure		500				{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/perda_pendirian [get]
func (h *PerdaPendirianHandler) Index(c *fiber.Ctx) error {
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
	modalDasarMin := c.QueryFloat("modal_dasar_min")
	modalDasarMax := c.QueryFloat("modal_dasar_max")

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idBumd,
		page,
		limit,
		search,
		modalDasarMin,
		modalDasarMax,
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

// View func for get data perda pendirian by id.
//
//	@Summary		get data perda pendirian by id
//	@Description	get data perda pendirian by id.
//	@ID				perda_pendirian-view
//	@Tags			Perda Pendirian
//	@Produce		json
//	@Param			id_bumd	path		int							true	"Id BUMD"
//	@Param			id		path		int							true	"Id untuk get data perda pendirian"
//	@success		200		{object}	dokumen.PerdaPendirianModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/perda_pendirian/{id} [get]
func (h *PerdaPendirianHandler) View(c *fiber.Ctx) error {
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

// Create func for create data perda pendirian.
//
//	@Summary		create data perda pendirian
//	@Description	create data perda pendirian.
//	@ID				perda_pendirian-create
//	@Tags			Perda Pendirian
//	@Accept			multipart/form-data
//	@Param			id_bumd			path		int		true	"Id BUMD"
//	@Param			nomor_perda		formData	string	true	"Nomor Perda"
//	@Param			tanggal_perda	formData	string	true	"Tanggal Perda"
//	@Param			keterangan		formData	string	true	"Keterangan"
//	@Param			file			formData	file	false	"File"
//	@Param			modal_dasar		formData	string	true	"Modal Dasar"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/perda_pendirian [post]
func (h *PerdaPendirianHandler) Create(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	payload := new(dokumen.PerdaPendirianForm)
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

// Update func for update data perda pendirian.
//
//	@Summary		update data perda pendirian
//	@Description	update data perda pendirian.
//	@ID				perda_pendirian-update
//	@Tags			Perda Pendirian
//	@Accept			multipart/form-data
//	@Param			id_bumd			path		int		true	"Id BUMD"
//	@Param			id				path		int		true	"Id untuk update data perda pendirian"
//	@Param			nomor_perda		formData	string	true	"Nomor Perda"
//	@Param			tanggal_perda	formData	string	true	"Tanggal Perda"
//	@Param			keterangan		formData	string	true	"Keterangan"
//	@Param			file			formData	file	false	"File"
//	@Param			modal_dasar		formData	string	true	"Modal Dasar"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/perda_pendirian [put]
func (h *PerdaPendirianHandler) Update(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	payload := new(dokumen.PerdaPendirianForm)
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

// Delete func for delete data perda pendirian.
//
//	@Summary		delete data perda pendirian
//	@Description	delete data perda pendirian.
//	@ID				perda_pendirian-delete
//	@Tags			Perda Pendirian
//	@Accept			json
//	@Param			id_bumd	path	int	true	"Id BUMD"
//	@Param			id		path	int	true	"Id untuk delete data perda pendirian"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/perda_pendirian/{id} [delete]
func (h *PerdaPendirianHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	idBumd, err := c.ParamsInt("id_bumd")
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
