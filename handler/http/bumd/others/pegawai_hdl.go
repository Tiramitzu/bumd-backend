package others

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/others"
	"microdata/kemendagri/bumd/models/bumd/others"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type PegawaiHandler struct {
	Controller *ctl.PegawaiController
	Validate   *validator.Validate
}

func NewPegawaiHandler(app *fiber.App, controller *ctl.PegawaiController, vld *validator.Validate) {
	handler := &PegawaiHandler{
		Controller: controller,
		Validate:   vld,
	}

	rPub := app.Group("/pegawai")
	rPub.Get("/", handler.Index)
	rPub.Get("/:id", handler.View)
	rPub.Post("/", handler.Create)
	rPub.Put("/:id", handler.Update)
	rPub.Delete("/:id", handler.Delete)
}

// Index func for get all pegawai.
//
//	@Summary		get all pegawai
//	@Description	get all pegawai.
//	@ID				pegawai-index
//	@Tags			Pegawai
//	@Produce		json
//	@Param			id_bumd	query		int						true	"Id Bumd"
//	@Param			page	query		int						false	"Halaman yang ditampilkan"
//	@Param			limit	query		int						false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string					false	"Search"
//	@success		200		{object}	[]others.PegawaiModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pegawai [get]
func (h *PegawaiHandler) Index(c *fiber.Ctx) error {
	idBumd := c.QueryInt("id_bumd", 0)
	page := c.QueryInt("page", 1)
	search := c.Query("search", "")
	var limit int
	limit = c.QueryInt("limit", 5)

	if limit > 5 {
		limit = 5
	}

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		page,
		limit,
		idBumd,
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

// View func for get pegawai by id.
//
//	@Summary		get pegawai by id
//	@Description	get pegawai by id.
//	@ID				pegawai-view
//	@Tags			Pegawai
//	@Produce		json
//	@Param			id	path		int						true	"Id Pegawai"
//	@success		200	{object}	others.PegawaiModel	"Success"
//	@Failure		400	{object}	utils.RequestError		"Bad request"
//	@Failure		404	{object}	utils.RequestError		"Data not found"
//	@Failure		422	{array}		utils.RequestError		"Data validation failed"
//	@Failure		500	{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pegawai/:id [get]
func (h *PegawaiHandler) View(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}

	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create pegawai.
//
//	@Summary		create pegawai
//	@Description	create pegawai.
//	@ID				pegawai-create
//	@Tags			Pegawai
//	@Accept			json
//	@Param			payload	body	others.PegawaiForm	true	"Pegawai payload"
//	@success		200						{object}	others.PegawaiModel	"Success"
//	@Failure		400						{object}	utils.RequestError		"Bad request"
//	@Failure		404						{object}	utils.RequestError		"Data not found"
//	@Failure		422						{array}		utils.RequestError		"Data validation failed"
//	@Failure		500						{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pegawai [post]
func (h *PegawaiHandler) Create(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}

	payload := new(others.PegawaiForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update pegawai.
//
//	@Summary		update pegawai
//	@Description	update pegawai.
//	@ID				pegawai-update
//	@Tags			Pegawai
//	@Accept			json
//	@Param			id_bumd					path		int						true	"Id Bumd"
//	@Param			payload	body	others.PegawaiForm	true	"Pegawai payload"
//	@success		200						{object}	others.PegawaiModel	"Success"
//	@Failure		400						{object}	utils.RequestError		"Bad request"
//	@Failure		404						{object}	utils.RequestError		"Data not found"
//	@Failure		422						{array}		utils.RequestError		"Data validation failed"
//	@Failure		500						{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pegawai/:id [put]
func (h *PegawaiHandler) Update(c *fiber.Ctx) error {
	idBumd, err := c.ParamsInt("id_bumd")
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(others.PegawaiForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, id, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete pegawai.
//
//	@Summary		delete pegawai
//	@Description	delete pegawai.
//	@ID				pegawai-delete
//	@Tags			Pegawai
//	@Param			id_bumd	path		int						true	"Id Bumd"
//	@Param			id		path		int						true	"Id Pegawai"
//	@success		200		{object}	others.PegawaiModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pegawai/:id [delete]
func (h *PegawaiHandler) Delete(c *fiber.Ctx) error {
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
