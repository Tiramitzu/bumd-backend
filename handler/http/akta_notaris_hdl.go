package http

import (
	"microdata/kemendagri/bumd/controller"
	"microdata/kemendagri/bumd/models"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AktaNotarisHandler struct {
	Controller *controller.AktaNotarisController
	Validate   *validator.Validate
}

func NewAktaNotarisHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *controller.AktaNotarisController,
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
//	@Param			id_bumd	query		int						true	"Id BUMD"
//	@Param			page	query		int						false	"Halaman yang ditampilkan"
//	@Param			limit	query		int						false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string					false	"Pencarian"
//	@success		200		{object}	models.AktaNotarisModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/akta_notaris [get]
func (h *AktaNotarisHandler) Index(c *fiber.Ctx) error {
	idBumd := c.QueryInt("id_bumd")
	if idBumd < 1 {
		return utils.RequestError{
			Code:    fiber.StatusBadRequest,
			Message: "Salah satu field tidak valid",
			Fields: []utils.DataValidationError{
				{
					Field:   "id_bumd",
					Message: "Id BUMD harus harus diisi",
				},
			},
		}
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
//	@Param			id	path		int						true	"Id untuk get data akta notaris"
//	@success		200	{object}	models.AktaNotarisModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/akta_notaris/{id} [get]
func (h *AktaNotarisHandler) View(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), id)
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
//	@Param			nomor_perda		formData	string	true	"Nomor Perda"
//	@Param			tanggal_perda	formData	string	true	"Tanggal Perda"
//	@Param			keterangan		formData	string	true	"Keterangan"
//	@Param			file			formData	string	false	"File"
//	@Param			modal_dasar		formData	float64	true	"Modal Dasar"
//	@Param			id_bumd			formData	int		true	"Id BUMD"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/akta_notaris [post]
func (h *AktaNotarisHandler) Create(c *fiber.Ctx) error {
	payload := new(models.AktaNotarisForm)
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
//	@Param			nomor_perda		formData	string					true	"Nomor Perda"
//	@Param			tanggal_perda	formData	string					true	"Tanggal Perda"
//	@Param			keterangan		formData	string					true	"Keterangan"
//	@Param			file			formData	string					false	"File"
//	@Param			modal_dasar		formData	float64					true	"Modal Dasar"
//	@Param			id_bumd			formData	int						true	"Id BUMD"
//	@Param			id				path		int						true	"Id untuk update data akta notaris"
//	@Param			payload			body		models.AktaNotarisForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/akta_notaris/{id} [put]
func (h *AktaNotarisHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(models.AktaNotarisForm)
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
		payload,
		id,
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
//	@Param			id	path	int	true	"Id untuk delete data akta notaris"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/akta_notaris/{id} [delete]
func (h *AktaNotarisHandler) Delete(c *fiber.Ctx) error {
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
