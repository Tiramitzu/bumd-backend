package http_mst

import (
	controller_mst "microdata/kemendagri/bumd/controller/master"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JenisLaporanHandler struct {
	Controller *controller_mst.JenisLaporanController
	Validate   *validator.Validate
}

func NewJenisLaporanHandler(r fiber.Router, validator *validator.Validate, controller *controller_mst.JenisLaporanController) {
	handler := &JenisLaporanHandler{
		Controller: controller,
		Validate:   validator,
	}

	rStrict := r.Group("jenis_laporan")
	rStrict.Get("/", handler.Index)
}

// Index func for index.
//
//	@Summary		get jenis laporan
//	@Description	get jenis laporan.
//	@ID				jenis_laporan-index
//	@Tags			JenisLaporan
//	@Accept			json
//	@Produce		json
//	@Param			bentuk_usaha	query		string						false	"Bentuk Usaha"	Format(uuid)
//	@Param			parent_id		query		int							false	"Parent ID"
//	@Success		200				{object}	[]models.JenisLaporanModel	"Success"
//	@Failure		400				{object}	utils.RequestError			"Bad request"
//	@Failure		403				{object}	utils.LoginError			"Login forbidden"
//	@Failure		404				{object}	utils.RequestError			"Data not found"
//	@Failure		422				{array}		utils.RequestError			"Data validation failed"
//	@Failure		500				{object}	utils.RequestError			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/jenis_laporan [get]
func (h *JenisLaporanHandler) Index(c *fiber.Ctx) error {
	bentukUsaha := c.Query("bentuk_usaha", "00000000-0000-0000-0000-000000000000")
	bentukUsahaUUID, err := uuid.Parse(bentukUsaha)
	if err != nil {
		return err
	}
	parentId := c.QueryInt("parent_id")

	r, err := h.Controller.Index(c.Context(), c.Locals("jwt").(*jwt.Token), bentukUsahaUUID, parentId)
	if err != nil {
		return err
	}

	return c.JSON(r)
}
