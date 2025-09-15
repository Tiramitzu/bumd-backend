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

type ProdukHandler struct {
	Controller *ctl.ProdukController
	Validate   *validator.Validate
}

func NewProdukHandler(r fiber.Router, validator *validator.Validate, controller *ctl.ProdukController) {
	handler := &ProdukHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("produk")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get data produk.
//
//	@Summary		get data produk
//	@Description	get data produk.
//	@ID				produk-index
//	@Tags			Produk
//	@Produce		json
//	@Param			id_bumd	path		string					true	"Id BUMD"	Format(uuid)
//	@Param			page	query		int						false	"Page"
//	@Param			limit	query		int						false	"Limit"
//	@Param			search	query		string					false	"Search"
//	@Success		200		{object}	[]others.ProdukModel	"Success"
//	@Failure		400		{object}	utils.RequestError		"Bad request"
//	@Failure		404		{object}	utils.RequestError		"Data not found"
//	@Failure		422		{array}		utils.RequestError		"Data validation failed"
//	@Failure		500		{object}	utils.RequestError		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/produk [get]
func (h *ProdukHandler) Index(c *fiber.Ctx) error {
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

// View func for get data produk.
//
//	@Summary		get data produk
//	@Description	get data produk.
//	@ID				produk-view
//	@Tags			Produk
//	@Produce		json
//	@Param			id_bumd	path		string				true	"Id BUMD"	Format(uuid)
//	@Param			id		path		string				true	"Id PRODUK"	Format(uuid)
//	@Success		200		{object}	others.ProdukModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/produk/{id} [get]
func (h *ProdukHandler) View(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	id := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	m, err := h.Controller.View(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Create func for create data produk.
//
//	@Summary		create data produk
//	@Description	create data produk.
//	@ID				produk-create
//	@Tags			Produk
//	@Accept			multipart/form-data
//	@Param			id_bumd		path		string				true	"Id BUMD"	Format(uuid)
//	@Param			nama_produk	formData	string				false	"Nama Produk"
//	@Param			deskripsi	formData	int					false	"Deskripsi Produk"
//	@Param			foto_produk	formData	file				false	"Foto Produk"
//	@Success		200			{object}	bool				"Success"
//	@Failure		400			{object}	utils.RequestError	"Bad request"
//	@Failure		422			{array}		utils.RequestError	"Data validation failed"
//	@Failure		500			{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/produk [post]
func (h *ProdukHandler) Create(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	payload := new(others.ProdukForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Create(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update func for update data produk.
//
//	@Summary		update data produk
//	@Description	update data produk.
//	@ID				produk-update
//	@Tags			Produk
//	@Accept			multipart/form-data
//	@Param			id_bumd		path		string				true	"Id BUMD"	Format(uuid)
//	@Param			nama_produk	formData	string				false	"Nama Produk"
//	@Param			deskripsi	formData	int					false	"Deskripsi Produk"
//	@Param			foto_produk	formData	file				false	"Foto Produk"
//	@Success		200			{object}	bool				"Success"
//	@Failure		400			{object}	utils.RequestError	"Bad request"
//	@Failure		422			{array}		utils.RequestError	"Data validation failed"
//	@Failure		500			{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/produk/{id} [put]
func (h *ProdukHandler) Update(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	id := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	payload := new(others.ProdukForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	m, err := h.Controller.Update(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, parsedId, payload)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete func for delete data produk.
//
//	@Summary		delete data produk
//	@Description	delete data produk.
//	@ID				produk-delete
//	@Tags			Produk
//	@Accept			json
//	@Param			id_bumd	path		string				true	"Id BUMD"	Format(uuid)
//	@Param			id		path		string				true	"Id PRODUK"	Format(uuid)
//	@Success		200		{object}	bool				"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/bumd/{id_bumd}/produk/{id} [delete]
func (h *ProdukHandler) Delete(c *fiber.Ctx) error {
	idBumdStr := c.Params("id_bumd")
	id := c.Params("id")
	idBumd, err := uuid.Parse(idBumdStr)
	if err != nil {
		return err
	}
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m, err := h.Controller.Delete(c.Context(), c.Locals("jwt").(*jwt.Token), idBumd, parsedId)
	if err != nil {
		return err
	}
	return c.JSON(m)
}
