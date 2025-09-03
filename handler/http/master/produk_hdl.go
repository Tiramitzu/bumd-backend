package http_mst

import (
	controller_mst "microdata/kemendagri/bumd/controller/master"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type ProdukHandler struct {
	Controller *controller_mst.ProdukController
	Validate   *validator.Validate
}

func NewProdukHandler(r fiber.Router, validator *validator.Validate, controller *controller_mst.ProdukController) {
	handler := &ProdukHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("produk")
	rStrict.Get("/", handler.Index)
}

// Index func for get data produk.
//
//	@Summary		get data produk
//	@Description	get data produk.
//	@ID				produk-index
//	@Tags			Produk
//	@Produce		json
//	@Param			id_bumd	query		int					false	"Id BUMD"
//	@Param			nama	query		string				false	"Nama Produk"
//	@Param			page	query		int					false	"Halaman yang ditampilkan"
//	@Param			limit	query		int					false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@success		200		{object}	models.ProdukModel	"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/produk [get]
func (h *ProdukHandler) Index(c *fiber.Ctx) error {
	idBumd := c.QueryInt("id_bumd", 0)
	nama := c.Query("nama")
	page := c.QueryInt("page", 1)
	var limit int
	limit = c.QueryInt("limit", 5)

	if limit > 5 {
		limit = 5
	}

	m, totalCount, pageCount, err := h.Controller.Index(
		c.Context(),
		c.Locals("jwt").(*jwt.Token),
		idBumd,
		nama,
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
		c.Append("x-pagination-previous-page", strconv.Itoa(page-1))
	}
	c.Append("x-pagination-current-page", strconv.Itoa(page))
	if page < pageCount {
		c.Append("x-pagination-next-page", strconv.Itoa(page+1))
	}
	return c.JSON(m)
}
