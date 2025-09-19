package kepengurusan_sdm

import (
	ctl "microdata/kemendagri/bumd/controller/bumd/kepengurusan_sdm"
	"microdata/kemendagri/bumd/models/bumd/kepengurusan_sdm"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type PengurusHandler struct {
	Controller *ctl.PengurusController
	Validate   *validator.Validate
}

func NewPengurusHandler(r fiber.Router, vld *validator.Validate, controller *ctl.PengurusController) {
	handler := &PengurusHandler{
		Controller: controller,
		Validate:   vld,
	}

	rPub := r.Group("/pengurus")
	rPub.Get("/", handler.Index)
	rPub.Get("/:id", handler.View)
	rPub.Post("/", handler.Create)
	rPub.Put("/:id", handler.Update)
	rPub.Delete("/:id", handler.Delete)
}

// Index func for get all pengurus.
//
//	@Summary		get all pengurus
//	@Description	get all pengurus.
//	@ID				pengurus-index
//	@Tags			Pengurus
//	@Produce		json
//	@Param			id_bumd	path		string								true	"Id Bumd"	Format(uuid)
//	@Param			page	query		int									false	"Halaman yang ditampilkan"
//	@Param			limit	query		int									false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@Param			search	query		string								false	"Search"
//	@success		200		{object}	[]kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400		{object}	utils.RequestError					"Bad request"
//	@Failure		404		{object}	utils.RequestError					"Data not found"
//	@Failure		422		{array}		utils.RequestError					"Data validation failed"
//	@Failure		500		{object}	utils.RequestError					"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus [get]
func (h *PengurusHandler) Index(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}
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
		parsedIdBumd,
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

// View func for get pengurus by id.
//
//	@Summary		get pengurus by id
//	@Description	get pengurus by id.
//	@ID				pengurus-view
//	@Tags			Pengurus
//	@Produce		json
//	@Param			id		path		string							true	"Id Pengurus"	Format(uuid)
//	@Param			id_bumd	path		string							true	"Id Bumd"		Format(uuid)
//	@success		200		{object}	kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		422		{array}		utils.RequestError				"Data validation failed"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus/{id} [get]
func (h *PengurusHandler) View(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}

	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create pengurus.
//
//	@Summary		create pengurus
//	@Description	create pengurus.
//	@ID				pengurus-create
//	@Tags			Pengurus
//	@Accept			multipart/form-data
//	@Param			id_bumd					path		string							true	"Id Bumd"	Format(uuid)
//	@Param			jabatan_struktur		formData	int								true	"Jabatan Struktur"
//	@Param			nama_pengurus			formData	string							true	"Nama Pengurus"
//	@Param			nik						formData	string							true	"NIK"
//	@Param			alamat					formData	string							true	"Alamat"
//	@Param			deskripsi_jabatan		formData	string							true	"Deskripsi Jabatan"
//	@Param			pendidikan_akhir		formData	string							true	"Pendidikan Akhir"	Format(uuid)
//	@Param			tanggal_mulai_jabatan	formData	string							true	"Tanggal Mulai Jabatan"
//	@Param			tanggal_akhir_jabatan	formData	string							true	"Tanggal Akhir Jabatan"
//	@Param			file					formData	file							true	"File"
//	@success		200						{object}	kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400						{object}	utils.RequestError				"Bad request"
//	@Failure		404						{object}	utils.RequestError				"Data not found"
//	@Failure		422						{array}		utils.RequestError				"Data validation failed"
//	@Failure		500						{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus [post]
func (h *PengurusHandler) Create(c *fiber.Ctx) error {
	idBumd := c.Params("id_bumd")
	parsedIdBumd, err := uuid.Parse(idBumd)
	if err != nil {
		return err
	}

	payload := new(kepengurusan_sdm.PengurusForm)

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
		if err := utils.ValidateFile(payload.File, maxFileSize, []string{"application/pdf", "image/jpeg", "image/png", "image/jpg"}); err != nil {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update pengurus.
//
//	@Summary		update pengurus
//	@Description	update pengurus.
//	@ID				pengurus-update
//	@Tags			Pengurus
//	@Accept			multipart/form-data
//	@Param			id_bumd					path		string							true	"Id Bumd"			Format(uuid)
//	@Param			id						path		string							true	"Id Pengurus"		Format(uuid)
//	@Param			jabatan_struktur		formData	int								true	"Jabatan Struktur"	Min(0)	Max(2)
//	@Param			nama_pengurus			formData	string							true	"Nama Pengurus"
//	@Param			nik						formData	string							true	"NIK"
//	@Param			alamat					formData	string							true	"Alamat"
//	@Param			deskripsi_jabatan		formData	string							true	"Deskripsi Jabatan"
//	@Param			pendidikan_akhir		formData	string							true	"Pendidikan Akhir"	Format(uuid)
//	@Param			tanggal_mulai_jabatan	formData	string							true	"Tanggal Mulai Jabatan"
//	@Param			tanggal_akhir_jabatan	formData	string							true	"Tanggal Akhir Jabatan"
//	@Param			file					formData	file							false	"File"
//	@success		200						{object}	kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400						{object}	utils.RequestError				"Bad request"
//	@Failure		404						{object}	utils.RequestError				"Data not found"
//	@Failure		422						{array}		utils.RequestError				"Data validation failed"
//	@Failure		500						{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus/{id} [put]
func (h *PengurusHandler) Update(c *fiber.Ctx) error {
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

	formModel := new(kepengurusan_sdm.PengurusForm)

	// Parse form fields first
	if err := c.BodyParser(formModel); err != nil {
		return err
	}

	// Handle file upload separately
	file, err := c.FormFile("file")
	if err == nil {
		// If file is found, assign it to payload
		formModel.File = file
	}

	if err := h.Validate.Struct(formModel); err != nil {
		return err
	}

	if formModel.File != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(formModel.File, maxFileSize, []string{"application/pdf", "image/jpeg", "image/png", "image/jpg"}); err != nil {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			}
		}
	}

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId, formModel)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete pengurus.
//
//	@Summary		delete pengurus
//	@Description	delete pengurus.
//	@ID				pengurus-delete
//	@Tags			Pengurus
//	@Param			id_bumd	path		string							true	"Id Bumd"		Format(uuid)
//	@Param			id		path		string							true	"Id Pengurus"	Format(uuid)
//	@success		200		{object}	kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		422		{array}		utils.RequestError				"Data validation failed"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus/{id} [delete]
func (h *PengurusHandler) Delete(c *fiber.Ctx) error {
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

// UpdateStatus func for update status pengurus.
//
//	@Summary		update status pengurus
//	@Description	update status pengurus.
//	@ID				pengurus-update-status
//	@Tags			Pengurus
//	@Param			id_bumd	path		string							true	"Id Bumd"		Format(uuid)
//	@Param			id		path		string							true	"Id Pengurus"	Format(uuid)
//	@success		200		{object}	kepengurusan_sdm.PengurusModel	"Success"
//	@Failure		400		{object}	utils.RequestError				"Bad request"
//	@Failure		404		{object}	utils.RequestError				"Data not found"
//	@Failure		422		{array}		utils.RequestError				"Data validation failed"
//	@Failure		500		{object}	utils.RequestError				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/pengurus/{id}/status [put]
func (h *PengurusHandler) UpdateStatus(c *fiber.Ctx) error {
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

	payload := new(kepengurusan_sdm.PengurusUpdateForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	m, err := h.Controller.UpdateStatus(c.Context(), c.Locals("jwt").(*jwt.Token), parsedIdBumd, parsedId, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
