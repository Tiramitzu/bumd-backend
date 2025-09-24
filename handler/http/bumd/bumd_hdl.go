package bumd

import (
	ctl "microdata/kemendagri/bumd/controller/bumd"
	ctl_dkmn "microdata/kemendagri/bumd/controller/bumd/dokumen"
	ctl_kepengurusan_sdm "microdata/kemendagri/bumd/controller/bumd/kepengurusan_sdm"
	ctl_keu "microdata/kemendagri/bumd/controller/bumd/keuangan"
	ctl_others "microdata/kemendagri/bumd/controller/bumd/others"
	"microdata/kemendagri/bumd/handler/http/bumd/dokumen"
	"microdata/kemendagri/bumd/handler/http/bumd/kepengurusan_sdm"
	"microdata/kemendagri/bumd/handler/http/bumd/keuangan"
	"microdata/kemendagri/bumd/handler/http/bumd/others"
	"microdata/kemendagri/bumd/models/bumd"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BumdHandler struct {
	Controller *ctl.BumdController
	Validate   *validator.Validate
	pgxConn    *pgxpool.Pool
	minioConn  *utils.MinioConn
}

func NewBumdHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *ctl.BumdController,
	pgxConn *pgxpool.Pool,
	minioConn *utils.MinioConn,
) {
	handler := &BumdHandler{
		Controller: controller,
		Validate:   validator,
		pgxConn:    pgxConn,
		minioConn:  minioConn,
	}

	rStrict := r.Group("bumd")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/kelengkapan_input", handler.KelengkapanInput)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)

	rStrict.Get("/:id/logo", handler.Logo)
	rStrict.Put("/:id/logo", handler.LogoUpdate)

	rStrict.Get("/:id/spi", handler.SPI)
	rStrict.Put("/:id/spi", handler.SPIUpdate)

	rStrict.Get("/:id/npwp", handler.NPWP)
	rStrict.Put("/:id/npwp", handler.NPWPUpdate)

	rData := rStrict.Group("/:id_bumd")
	dokumen.NewPerdaPendirianHandler(
		rData,
		validator,
		ctl_dkmn.NewPerdaPendirianController(pgxConn, minioConn),
	)
	dokumen.NewAktaNotarisHandler(
		rData,
		validator,
		ctl_dkmn.NewAktaNotarisController(pgxConn, minioConn),
	)
	dokumen.NewSiupHandler(
		rData,
		validator,
		ctl_dkmn.NewSiupController(pgxConn, minioConn),
	)
	dokumen.NewNibHandler(
		rData,
		validator,
		ctl_dkmn.NewNibController(pgxConn, minioConn),
	)
	dokumen.NewTdpHandler(
		rData,
		validator,
		ctl_dkmn.NewTdpController(pgxConn, minioConn),
	)

	// keuangan
	keuangan.NewModalHandler(
		rData,
		validator,
		ctl_keu.NewModalController(pgxConn),
	)
	keuangan.NewLaporanKeuanganHandler(
		rData,
		validator,
		ctl_keu.NewLaporanKeuanganController(pgxConn, minioConn),
	)

	// others
	others.NewDomisiliHandler(
		rData,
		validator,
		ctl_others.NewDomisiliController(pgxConn, minioConn),
	)
	others.NewRKAHandler(
		rData,
		validator,
		ctl_others.NewRKAController(pgxConn, minioConn),
	)
	others.NewRencanaBisnisHandler(
		rData,
		validator,
		ctl_others.NewRencanaBisnisController(pgxConn, minioConn),
	)
	others.NewPeraturanHandler(
		rData,
		validator,
		ctl_others.NewPeraturanController(pgxConn, minioConn),
	)
	others.NewProdukHandler(
		rData,
		validator,
		ctl_others.NewProdukController(pgxConn, minioConn),
	)
	others.NewKinerjaHandler(
		rData,
		validator,
		ctl_others.NewKinerjaController(pgxConn),
	)

	// kepengurusan sdm
	kepengurusan_sdm.NewPengurusHandler(
		rData,
		validator,
		ctl_kepengurusan_sdm.NewPengurusController(pgxConn, minioConn),
	)
	kepengurusan_sdm.NewPegawaiHandler(
		rData,
		validator,
		ctl_kepengurusan_sdm.NewPegawaiController(pgxConn),
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
//	@Param			induk_perusahaan	query		string				false	"Induk Perusahaan"	Format(uuid)
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
	indukPerusahaan := c.Query("induk_perusahaan", uuid.Nil.String())
	parsedIndukPerusahaan, err := uuid.Parse(indukPerusahaan)
	if err != nil {
		return err
	}

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
		parsedIndukPerusahaan,
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
//	@Param			id	path		string				true	"Id untuk get data bumd"	Format(uuid)
//	@success		200	{object}	bumd.BumdModel		"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [get]
func (h *BumdHandler) View(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), parsedId)
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
//	@Param			id		path	string			true	"Id untuk update data bumd"	Format(uuid)
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
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
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
		parsedId,
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
//	@Param			id	path	string	true	"Id untuk delete data bumd"	Format(uuid)
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id} [delete]
func (h *BumdHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.Delete(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		parsedId,
	)
	if err != nil {
		return err
	}

	return c.JSON(m)
}

// Logo func for get data logo by id.
//
//	@Summary		get data logo by id
//	@Description	get data logo by id.
//	@ID				logo-view
//	@Tags			BUMD
//	@Produce		json
//	@Param			id	path		string				true	"Id untuk get data logo"	Format(uuid)
//	@success		200	{object}	bumd.LogoModel		"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/logo [get]
func (h *BumdHandler) Logo(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.Logo(c.Context(), c.Locals("jwt").(*jwt.Token), parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// LogoUpdate func for update data logo by id.
//
//	@Summary		update data logo by id
//	@Description	update data logo by id.
//	@ID				logo-update
//	@Tags			BUMD
//	@Accept			multipart/form-data
//	@Param			id		path		string	true	"Id untuk update data logo"	Format(uuid)
//	@Param			file	formData	file	false	"File"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/logo [put]
func (h *BumdHandler) LogoUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	payload := new(bumd.LogoForm)

	// Parse form fields first
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// Handle file upload separately
	file, err := c.FormFile("file")
	if err == nil {
		// If file is found, assign it to payload
		payload.Logo = file
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.LogoUpdate(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
		parsedId,
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// SPI func for get data spi by id.
//
//	@Summary		get data spi by id
//	@Description	get data spi by id.
//	@ID				spi-view
//	@Tags			BUMD
//	@Produce		json
//	@Param			id	path		string				true	"Id untuk get data spi"	Format(uuid)
//	@success		200	{object}	bumd.SPIModel		"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/spi [get]
func (h *BumdHandler) SPI(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.SPI(c.Context(), c.Locals("jwt").(*jwt.Token), parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// SPIUpdate func for update data spi by id.
//
//	@Summary		update data spi by id
//	@Description	update data spi by id.
//	@ID				spi-update
//	@Tags			BUMD
//	@Accept			multipart/form-data
//	@Param			id				path		string	true	"Id untuk update data spi"	Format(uuid)
//	@Param			penerapan_spi	formData	bool	false	"Penerapan SPI"
//	@Param			file_spi		formData	file	false	"File SPI"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/spi [put]
func (h *BumdHandler) SPIUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	payload := new(bumd.SPIForm)

	// Parse form fields first
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// Handle file upload separately
	file, err := c.FormFile("file_spi")
	if err == nil {
		// If file is found, assign it to payload
		payload.FileSPI = file
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.FileSPI != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.FileSPI, maxFileSize, []string{"application/pdf"}); err != nil {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	m, err := h.Controller.SPIUpdate(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
		parsedId,
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// NPWP func for get data npwp by id.
//
//	@Summary		get data npwp by id
//	@Description	get data npwp by id.
//	@ID				npwp-view
//	@Tags			BUMD
//	@Produce		json
//	@Param			id	path		string				true	"Id untuk get data npwp"	Format(uuid)
//	@success		200	{object}	bumd.NPWPModel		"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/npwp [get]
func (h *BumdHandler) NPWP(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.NPWP(c.Context(), c.Locals("jwt").(*jwt.Token), parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// NPWPUpdate func for update data npwp by id.
//
//	@Summary		update data npwp by id
//	@Description	update data npwp by id.
//	@ID				npwp-update
//	@Tags			BUMD
//	@Accept			json
//	@Param			id		path		string	true	"Id bumd untuk update data npwp"	Format(uuid)
//	@Param			npwp	formData	string	false	"NPWP"
//	@Param			pemberi	formData	string	false	"Pemberi"
//	@Param			file	formData	file	false	"File"
//	@Produce		json
//	@success		200	{object}	boolean				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id}/npwp [put]
func (h *BumdHandler) NPWPUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	payload := new(bumd.NPWPForm)

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

	m, err := h.Controller.NPWPUpdate(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		payload,
		parsedId,
	)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// KelengkapanInput func for get data kelengkapan input.
//
//	@Summary		get data kelengkapan input
//	@Description	get data kelengkapan input.
//	@ID				kelengkapan_input-view
//	@Tags			BUMD
//	@Produce		json
//	@Param			id_daerah		query		string						false	"Id Daerah"
//	@Param			bentuk_usaha	query		string						false	"Bentuk Usaha BUMD"
//	@Param			id_bumd			query		string						false	"Id BUMD"
//	@Param			page			query		int							false	"Page"
//	@Param			limit			query		int							false	"Limit"
//	@Param			search			query		string						false	"Search"
//	@success		200				{object}	bumd.KelengkapanInputModel	"Success"
//	@Failure		400				{object}	utils.RequestError			"Bad request"
//	@Failure		404				{object}	utils.RequestError			"Data not found"
//	@Failure		500				{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/kelengkapan_input [get]
func (h *BumdHandler) KelengkapanInput(c *fiber.Ctx) error {
	idDaerah := c.QueryInt("id_daerah")
	bentukUsaha := c.Query("bentuk_usaha", uuid.Nil.String())
	parsedBentukUsaha, err := uuid.Parse(bentukUsaha)
	if err != nil {
		return err
	}
	idBumd := c.Query("id_bumd", uuid.Nil.String())
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
	page := c.QueryInt("page", 1)
	var limit int
	limit = c.QueryInt("limit", 5)
	search := c.Query("search", "")

	if limit > 5 {
		limit = 5
	}

	m, totalCount, pageCount, err := h.Controller.KelengkapanInput(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idDaerah,
		parsedBentukUsaha,
		parsedIdBumd,
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
