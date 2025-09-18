package others

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/others"
	"microdata/kemendagri/bumd/models/bumd/others"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type KinerjaHandler struct {
	Controller *ctl.KinerjaController
	Validate   *validator.Validate
}

func NewKinerjaHandler(r fiber.Router, validator *validator.Validate, controller *ctl.KinerjaController) {
	handler := &KinerjaHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("kinerja")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data kinerja.
//
//	@Summary		get data kinerja
//	@Description	get data kinerja.
//	@ID				kinerja-index
//	@Tags			Kinerja
//	@Produce		json
//	@Param			id_bumd	path		string					true	"Id BUMD"	Format(uuid)
//	@Param			tahun	query		int						false	"Tahun"
//	@Param			page	query		int						false	"Halaman yang ditampilkan"
//	@Param			limit	query		int						false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@success		200		{object}	[]others.KinerjaModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/kinerja [get]
func (h *KinerjaHandler) Index(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	page := c.QueryInt("page", 1)
	var limit int
	limit = c.QueryInt("limit", 5)
	if limit > 5 {
		limit = 5
	}
	tahun := c.QueryInt("tahun", 0)

	r, totalCount, pageCount, err := h.Controller.Index(c.Context(), c.Locals("jwt").(*jwt.Token), page, limit, idBumd, tahun)
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

	return c.JSON(r)
}

// View func for get data kinerja.
//
//	@Summary		get data kinerja
//	@Description	get data kinerja.
//	@ID				kinerja-view
//	@Tags			Kinerja
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"		Format(uuid)
//	@Param			id		path		string				true	"Id Kinerja"	Format(uuid)
//	@success		200		{object}	others.KinerjaModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/kinerja/{id} [get]
func (h *KinerjaHandler) View(c *fiber.Ctx) error {
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

	r, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}

	return c.JSON(r)
}

// Create func for create data kinerja.
//
//	@Summary		create data kinerja
//	@Description	create data kinerja.
//	@ID				kinerja-create
//	@Tags			Kinerja
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"	Format(uuid)
//	@Param			payload	body		others.KinerjaForm	true	"Payload"
//	@success		200		{object}	others.KinerjaModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/kinerja [post]
func (h *KinerjaHandler) Create(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	payload := new(others.KinerjaForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	r, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, payload)
	if err != nil {
		return err
	}

	return c.JSON(r)
}

// Update func for update data kinerja.
//
//	@Summary		update data kinerja
//	@Description	update data kinerja.
//	@ID				kinerja-update
//	@Tags			Kinerja
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"		Format(uuid)
//	@Param			id		path		string				true	"Id Kinerja"	Format(uuid)
//	@Param			payload	body		others.KinerjaForm	true	"Payload"
//	@success		200		{object}	others.KinerjaModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/kinerja/{id} [put]
func (h *KinerjaHandler) Update(c *fiber.Ctx) error {
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
	payload := new(others.KinerjaForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	r, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id, payload)
	if err != nil {
		return err
	}

	return c.JSON(r)
}

// Delete func for delete data kinerja.
//
//	@Summary		delete data kinerja
//	@Description	delete data kinerja.
//	@ID				kinerja-delete
//	@Tags			Kinerja
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"		Format(uuid)
//	@Param			id		path		string				true	"Id Kinerja"	Format(uuid)
//	@success		200		{object}	others.KinerjaModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/kinerja/{id} [delete]
func (h *KinerjaHandler) Delete(c *fiber.Ctx) error {
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

	r, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}

	return c.JSON(r)
}
