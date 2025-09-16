package others

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/others"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type RencanaBisnisHandler struct {
	Controller *ctl.RencanaBisnisController
	Validate   *validator.Validate
}

func NewRencanaBisnisHandler(r fiber.Router, validator *validator.Validate, controller *ctl.RencanaBisnisController) {
	handler := &RencanaBisnisHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("rencana_bisnis")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data rencana bisnis.
//
//	@Summary		get data rencana bisnis
//	@Description	get data rencana bisnis.
//	@ID				rencana_bisnis-index
//	@Tags			Rencana Bisnis
//	@Produce		json
//	@Param			id_bumd			path		string						true	"Id BUMD"	Format(uuid)
//	@Param			page			query		int							false	"Page"
//	@Param			limit			query		int							false	"Limit"
//	@Param			search			query		string						false	"Search"
//	@Param			kualifikasi		query		int							false	"Kualifikasi, 1 = Kecil, 2 = Non Kecil"
//	@Param			is_seumur_hidup	query		int							false	"Is Seumur Hidup"
//	@Success		200				{object}	[]others.RencanaBisnisModel	"Success"
//	@Failure		400				{object}	utils.RequestError			"Bad request"
//	@Failure		404				{object}	utils.RequestError			"Data not found"
//	@Failure		422				{array}		utils.RequestError			"Data validation failed"
//	@Failure		500				{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/rencana_bisnis [get]
func (h *RencanaBisnisHandler) Index(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idBumd, err := uuid.Parse(idBumdStr)
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

// View func for get data rencana bisnis.
//
//	@Summary		get data rencana bisnis
//	@Description	get data rencana bisnis.
//	@ID				rencana_bisnis-view
//	@Tags			Rencana Bisnis
//	@Produce		json
//	@Param			id_bumd	path		string						true	"Id BUMD"			Format(uuid)
//	@Param			id		path		string						true	"Id RENCANA BISNIS"	Format(uuid)
//	@Success		200		{object}	others.RencanaBisnisModel	"Success"
//	@Failure		400		{object}	utils.RequestError			"Bad request"
//	@Failure		404		{object}	utils.RequestError			"Data not found"
//	@Failure		422		{array}		utils.RequestError			"Data validation failed"
//	@Failure		500		{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/rencana_bisnis/{id} [get]
func (h *RencanaBisnisHandler) View(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idStr := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create data domisili.
//
//	@Summary		create data rencana bisnis
//	@Description	create data rencana bisnis.
//	@ID				rencana_bisnis-create
//	@Tags			Rencana Bisnis
//	@Accept			multipart/form-data
//	@Param			id_bumd				path		string				true	"Id BUMD"	Format(uuid)
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
//	@Router			/strict/bumd/{id_bumd}/rencana_bisnis [post]
func (h *RencanaBisnisHandler) Create(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	payload := new(others.RencanaBisnisForm)

	// Parse form fields first
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// Handle file upload separately
	file, err := c.FormFile("file")
	if err == nil {
		// If file is found, assign it to payload
		payload.File = file
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

// Update func for update data domisili.
//
//	@Summary		update data rencana bisnis
//	@Description	update data rencana bisnis.
//	@ID				rencana_bisnis-update
//	@Tags			Rencana Bisnis
//	@Accept			multipart/form-data
//	@Param			id_bumd				path		string				true	"Id BUMD"			Format(uuid)
//	@Param			id					path		string				true	"Id RENCANA BISNIS"	Format(uuid)
//	@Param			nomor				formData	string				true	"Nomor"
//	@Param			instansi_pemberi	formData	string				true	"Instansi Pemberi"
//	@Param			tanggal				formData	string				true	"Tanggal"
//	@Param			klasifikasi			formData	string				true	"Klasifikasi"
//	@Param			kualifikasi			formData	int					true	"Kualifikasi"
//	@Param			masa_berlaku		formData	string				false	"Masa Berlaku"
//	@Param			file				formData	file				false	"File"
//	@Success		200					{object}	bool				"Success"
//	@Failure		400					{object}	utils.RequestError	"Bad request"
//	@Failure		422					{array}		utils.RequestError	"Data validation failed"
//	@Failure		500					{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/rencana_bisnis/{id} [put]
func (h *RencanaBisnisHandler) Update(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idStr := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	payload := new(others.RencanaBisnisForm)

	// Parse form fields first
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// Handle file upload separately
	file, err := c.FormFile("file")
	if err == nil {
		// If file is found, assign it to payload
		payload.File = file
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

// Delete func for delete data domisili.
//
//	@Summary		delete data rencana bisnis
//	@Description	delete data rencana bisnis.
//	@ID				rencana_bisnis-delete
//	@Tags			Rencana Bisnis
//	@Accept			json
//	@Param			id_bumd	path		string				true	"Id BUMD"			Format(uuid)
//	@Param			id		path		string				true	"Id RENCANA BISNIS"	Format(uuid)
//	@Success		200		{object}	bool				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/rencana_bisnis/{id} [delete]
func (h *RencanaBisnisHandler) Delete(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idStr := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	m, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
