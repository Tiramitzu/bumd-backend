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

type LaporanKeuanganHandler struct {
	Controller *ctl.LaporanKeuanganController
	Validate   *validator.Validate
}

func NewLaporanKeuanganHandler(r fiber.Router, validator *validator.Validate, controller *ctl.LaporanKeuanganController) {
	handler := &LaporanKeuanganHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("laporan_keuangan")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data laporan keuangan.
//
//	@Summary		get data laporan keuangan
//	@Description	get data laporan keuangan.
//	@ID				laporan_keuangan-index
//	@Tags			Laporan Keuangan
//	@Produce		json
//	@Param			id_bumd	path		string							true	"Id BUMD"	Format(uuid)
//	@Param			page	query		int								false	"Halaman yang ditampilkan"
//	@Param			limit	query		int								false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string							false	"Pencarian"
//	@success		200		{object}	keuangan.LaporanKeuanganModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/laporan_keuangan [get]
func (h *LaporanKeuanganHandler) Index(c *fiber.Ctx) error {
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
	// search := c.Query("search")

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedIdBumd,
		page,
		limit,
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

// View func for get data laporan keuangan by id.
//
//	@Summary		get data laporan keuangan by id
//	@Description	get data laporan keuangan by id.
//	@ID				laporan_keuangan-view
//	@Tags			Laporan Keuangan
//	@Produce		json
//	@Param			id_bumd	path		string							true	"Id BUMD"				Format(uuid)
//	@Param			id		path		string							true	"Id Laporan Keuangan"	Format(uuid)
//	@success		200		{object}	keuangan.LaporanKeuanganModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/laporan_keuangan/{id} [get]
func (h *LaporanKeuanganHandler) View(c *fiber.Ctx) error {
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

// Create func for create data laporan keuangan.
//
//	@Summary		create data laporan keuangan
//	@Description	create data laporan keuangan.
//	@ID				laporan_keuangan-create
//	@Tags			Laporan Keuangan
//	@Accept			json
//	@Param			id_bumd	path	string							true	"Id BUMD"	Format(uuid)
//	@Param			payload	body	keuangan.LaporanKeuanganForm	true	"Create payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/laporan_keuangan [post]
func (h *LaporanKeuanganHandler) Create(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	payload := new(keuangan.LaporanKeuanganForm)

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

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update data laporan keuagan.
//
//	@Summary		update data laporan keuagan
//	@Description	update data laporan keuagan.
//	@ID				laporan_keuangan-update
//	@Tags			Laporan Keuangan
//	@Accept			json
//	@Param			id_bumd	path	string							true	"Id BUMD"	Format(uuid)
//	@Param			payload	body	keuangan.LaporanKeuanganForm	true	"Update payload"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/laporan_keuangan/{id} [put]
func (h *LaporanKeuanganHandler) Update(c *fiber.Ctx) error {
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
	payload := new(keuangan.LaporanKeuanganForm)

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

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete data laporan keuangan.
//
//	@Summary		delete data laporan keuangan
//	@Description	delete data laporan keuangan.
//	@ID				laporan_keuangan-delete
//	@Tags			Laporan Keuangan
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"				Format(uuid)
//	@Param			id		path		string				true	"Id Laporan Keuangan"	Format(uuid)
//	@success		200		{object}	boolean				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/laporan_keuangan/{id} [delete]
func (h *LaporanKeuanganHandler) Delete(c *fiber.Ctx) error {
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
