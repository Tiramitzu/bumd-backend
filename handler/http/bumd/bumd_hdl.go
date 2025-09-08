package bumd

import (
	ctl "microdata/kemendagri/bumd/controller/bumd"
	ctl_dkmn "microdata/kemendagri/bumd/controller/bumd/dokumen"
	ctl_keu "microdata/kemendagri/bumd/controller/bumd/keuangan"
	"microdata/kemendagri/bumd/handler/http/bumd/dokumen"
	"microdata/kemendagri/bumd/handler/http/bumd/keuangan"
	"microdata/kemendagri/bumd/models/bumd"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BumdHandler struct {
	Controller *ctl.BumdController
	Validate   *validator.Validate
	pgxConn    *pgxpool.Pool
}

func NewBumdHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *ctl.BumdController,
	pgxConn *pgxpool.Pool,
) {
	handler := &BumdHandler{
		Controller: controller,
		Validate:   validator,
		pgxConn:    pgxConn,
	}

	rStrict := r.Group("bumd")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)

	rData := rStrict.Group("/:id_bumd")
	dokumen.NewPerdaPendirianHandler(
		rData,
		validator,
		ctl_dkmn.NewPerdaPendirianController(pgxConn),
	)
	dokumen.NewAktaNotarisHandler(
		rData,
		validator,
		ctl_dkmn.NewAktaNotarisController(pgxConn),
	)
	dokumen.NewSiupHandler(
		rData,
		validator,
		ctl_dkmn.NewSiupController(pgxConn),
	)
	dokumen.NewNibHandler(
		rData,
		validator,
		ctl_dkmn.NewNibController(pgxConn),
	)
	dokumen.NewTdpHandler(
		rData,
		validator,
		ctl_dkmn.NewTdpController(pgxConn),
	)

	keuangan.NewModalHandler(
		rData,
		validator,
		ctl_keu.NewModalController(pgxConn),
	)
}

// Index func for get data bumd.
//
//	@Summary		get data bumd
//	@Description	get data bumd.
//	@ID				bumd-index
//	@Tags			BUMD
//	@Produce		json
//	@Param			nama				query		string				false	"Nama BUMD"
//	@Param			penerapan_spi		query		bool				false	"Penerapan SPI"
//	@Param			induk_perusahaan	query		int					false	"Induk Perusahaan"
//	@Param			page				query		int					false	"Halaman yang ditampilkan"
//	@Param			limit				query		int					false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@success		200					{object}	bumd.BumdModel		"Success"
//	@Failure		400					{object}	utils.RequestError	"Bad request"
//	@Failure		404					{object}	utils.RequestError	"Data not found"
//	@Failure		422					{array}		utils.RequestError	"Data validation failed"
//	@Failure		500					{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd [get]
func (h *BumdHandler) Index(c *fiber.Ctx) error {
	nama := c.Query("nama")
	penerapanSPI := c.QueryBool("penerapan_spi", false)
	indukPerusahaan := c.QueryInt("induk_perusahaan", 0)
	page := c.QueryInt("page", 1)
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
		nama,
		penerapanSPI,
		indukPerusahaan,
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

// View func for get data bumd by id.
//
//	@Summary		get data bumd by id
//	@Description	get data bumd by id.
//	@ID				bumd-view
//	@Tags			BUMD
//	@Produce		json
//	@Param			id	path		int					true	"Id untuk get data bumd"
//	@success		200	{object}	bumd.BumdModel		"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [get]
func (h *BumdHandler) View(c *fiber.Ctx) error {
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

// Create func for create data bumd.
//
//	@Summary		create data bumd
//	@Description	create data bumd.
//	@ID				bumd-create
//	@Tags			BUMD
//	@Accept			json
//	@Param			payload	body	bumd.BumdForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd [post]
func (h *BumdHandler) Create(c *fiber.Ctx) error {
	payload := new(bumd.BumdForm)
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

// Update func for update data bumd.
//
//	@Summary		update data bumd
//	@Description	update data bumd.
//	@ID				bumd-update
//	@Tags			BUMD
//	@Accept			json
//	@Param			id		path	int				true	"Id untuk update data bumd"
//	@Param			payload	body	bumd.BumdForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [put]
func (h *BumdHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(bumd.BumdForm)

	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
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

// Delete func for delete data bumd.
//
//	@Summary		delete data bumd
//	@Description	delete data bumd.
//	@ID				bumd-delete
//	@Tags			BUMD
//	@Accept			json
//	@Param			id	path	int	true	"Id untuk delete data bumd"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [delete]
func (h *BumdHandler) Delete(c *fiber.Ctx) error {
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
