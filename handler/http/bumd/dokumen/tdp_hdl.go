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

type TdpHandler struct {
	Controller *ctl.TdpController
	Validate   *validator.Validate
}

func NewTdpHandler(r fiber.Router, validator *validator.Validate, controller *ctl.TdpController) {
	handler := &TdpHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("tdp")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data tdp.
//
//	@Summary		get data tdp
//	@Description	get data tdp.
//	@ID				tdp-index
//	@Tags			Tdp
//	@Produce		json
//	@Param			id_bumd			path		int					true	"Id BUMD"
//	@Param			page			query		int					false	"Page"
//	@Param			limit			query		int					false	"Limit"
//	@Param			search			query		string				false	"Search"
//	@Param			kualifikasi		query		int					false	"Kualifikasi, 1 = Kecil, 2 = Non Kecil"
//	@Param			is_seumur_hidup	query		int					false	"Is Seumur Hidup"
//	@Success		200				{object}	[]dokumen.TdpModel	"Success"
//	@Failure		400				{object}	utils.RequestError	"Bad request"
//	@Failure		404				{object}	utils.RequestError	"Data not found"
//	@Failure		422				{array}		utils.RequestError	"Data validation failed"
//	@Failure		500				{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/tdp [get]
func (h *TdpHandler) Index(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	isSeumurHidup := c.QueryInt("is_seumur_hidup", 0)
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
		isSeumurHidup,
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

// View func for get data tdp.
//
//	@Summary		get data tdp
//	@Description	get data tdp.
//	@ID				tdp-view
//	@Tags			Tdp
//	@Produce		json
//	@Param			id_bumd	path		int					true	"Id BUMD"
//	@Param			id		path		int					true	"Id TDP"
//	@Success		200		{object}	dokumen.TdpModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/tdp/{id} [get]
func (h *TdpHandler) View(c *fiber.Ctx) error {
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

// Create func for create data tdp.
//
//	@Summary		create data tdp
//	@Description	create data tdp.
//	@ID				tdp-create
//	@Tags			Tdp
//	@Accept			multipart/form-data
//	@Param			id_bumd				path		int					true	"Id BUMD"
//	@Param			nomor				formData	string				false	"Nomor"
//	@Param			instansi_pemberi	formData	string				false	"Instansi Pemberi"
//	@Param			tanggal				formData	string				false	"Tanggal"
//	@Param			klasifikasi			formData	string				false	"Klasifikasi"
//	@Param			kualifikasi			formData	int					false	"Kualifikasi"
//	@Param			masa_berlaku		formData	string				false	"Masa Berlaku"
//	@Param			file				formData	file				false	"File"
//	@Success		200					{object}	bool				"Success"
//	@Failure		400					{object}	utils.RequestError	"Bad request"
//	@Failure		422					{array}		utils.RequestError	"Data validation failed"
//	@Failure		500					{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/tdp [post]
func (h *TdpHandler) Create(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	payload := new(dokumen.TdpForm)
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

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update data tdp.
//
//	@Summary		update data tdp
//	@Description	update data tdp.
//	@ID				tdp-update
//	@Tags			Tdp
//	@Accept			multipart/form-data
//	@Param			id_bumd				path		int					true	"Id BUMD"
//	@Param			id					path		int					true	"Id TDP"
//	@Param			nomor				formData	string				true	"Nomor"
//	@Param			instansi_pemberi	formData	string				true	"Instansi Pemberi"
//	@Param			tanggal				formData	string				true	"Tanggal"
//	@Param			klasifikasi			formData	string				true	"Klasifikasi"
//	@Param			kualifikasi			formData	int					true	"Kualifikasi"
//	@Param			masa_berlaku		formData	string				true	"Masa Berlaku"
//	@Param			file				formData	file				false	"File"
//	@Success		200					{object}	bool				"Success"
//	@Failure		400					{object}	utils.RequestError	"Bad request"
//	@Failure		422					{array}		utils.RequestError	"Data validation failed"
//	@Failure		500					{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/tdp/{id} [put]
func (h *TdpHandler) Update(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	payload := new(dokumen.TdpForm)
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

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete data tdp.
//
//	@Summary		delete data tdp
//	@Description	delete data tdp.
//	@ID				tdp-delete
//	@Tags			Tdp
//	@Accept			json
//	@Param			id_bumd	path		int					true	"Id BUMD"
//	@Param			id		path		int					true	"Id TDP"
//	@Success		200		{object}	bool				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/tdp/{id} [delete]
func (h *TdpHandler) Delete(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	m, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
