package others

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/others"
	"microdata/kemendagri/bumd/models/bumd/others"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type PeraturanHandler struct {
	Controller *ctl.PeraturanController
	Validate   *validator.Validate
}

func NewPeraturanHandler(r fiber.Router, validator *validator.Validate, controller *ctl.PeraturanController) {
	handler := &PeraturanHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("peraturan")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data peraturan.
//
//	@Summary		get data peraturan
//	@Description	get data peraturan.
//	@ID				peraturan-index
//	@Tags			Peraturan
//	@Produce		json
//	@Param			id_bumd	path		int						true	"Id BUMD"
//	@Param			page	query		int						false	"Page"
//	@Param			limit	query		int						false	"Limit"
//	@Param			search	query		string					false	"Search"
//	@Success		200		{object}	[]others.PeraturanModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/peraturan [get]
func (h *PeraturanHandler) Index(c *fiber.Ctx) error {
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

// View func for get data peraturan.
//
//	@Summary		get data peraturan
//	@Description	get data peraturan.
//	@ID				peraturan-view
//	@Tags			Peraturan
//	@Produce		json
//	@Param			id_bumd	path		int						true	"Id BUMD"
//	@Param			id		path		string					true	"Id PERATURAN"
//	@Success		200		{object}	others.PeraturanModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/peraturan/{id} [get]
func (h *PeraturanHandler) View(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id := c.Params("id")

	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create data peraturan.
//
//	@Summary		create data peraturan
//	@Description	create data peraturan.
//	@ID				peraturan-create
//	@Tags			Peraturan
//	@Accept			multipart/form-data
//	@Param			id_bumd					path		int					true	"Id BUMD"
//	@Param			nomor					formData	string				false	"Nomor"
//	@Param			jenis_peraturan			formData	int					false	"Jenis Peraturan"
//	@Param			tanggal_berlaku			formData	string				false	"Tanggal Berlaku"
//	@Param			keterangan_peraturan	formData	string				false	"Keterangan Peraturan"
//	@Param			file_peraturan			formData	file				false	"File Peraturan"
//	@Success		200						{object}	bool				"Success"
//	@Failure		400						{object}	utils.RequestError	"Bad request"
//	@Failure		422						{array}		utils.RequestError	"Data validation failed"
//	@Failure		500						{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/peraturan [post]
func (h *PeraturanHandler) Create(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	payload := new(others.PeraturanForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.FilePeraturan != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.FilePeraturan, maxFileSize, []string{"application/pdf"}); err != nil {
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

// Update func for update data peraturan.
//
//	@Summary		update data peraturan
//	@Description	update data peraturan.
//	@ID				peraturan-update
//	@Tags			Peraturan
//	@Accept			multipart/form-data
//	@Param			id_bumd					path		int					true	"Id BUMD"
//	@Param			id						path		string				true	"Id PERATURAN"
//	@Param			nomor					formData	string				true	"Nomor"
//	@Param			jenis_peraturan			formData	int					true	"Jenis Peraturan"
//	@Param			tanggal_berlaku			formData	string				false	"Tanggal Berlaku"
//	@Param			keterangan_peraturan	formData	string				true	"Keterangan Peraturan"
//	@Param			file_peraturan			formData	file				false	"File Peraturan"
//	@Success		200						{object}	bool				"Success"
//	@Failure		400						{object}	utils.RequestError	"Bad request"
//	@Failure		422						{array}		utils.RequestError	"Data validation failed"
//	@Failure		500						{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/peraturan/{id} [put]
func (h *PeraturanHandler) Update(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id := c.Params("id")
	payload := new(others.PeraturanForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.FilePeraturan != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.FilePeraturan, maxFileSize, []string{"application/pdf"}); err != nil {
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

// Delete func for delete data peraturan.
//
//	@Summary		delete data peraturan
//	@Description	delete data peraturan.
//	@ID				peraturan-delete
//	@Tags			Peraturan
//	@Accept			json
//	@Param			id_bumd	path		int					true	"Id BUMD"
//	@Param			id		path		string				true	"Id PERATURAN"
//	@Success		200		{object}	bool				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/peraturan/{id} [delete]
func (h *PeraturanHandler) Delete(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}
	id := c.Params("id")
	m, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
