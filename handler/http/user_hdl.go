package http

import (
	"fmt"
	"microdata/kemendagri/bumd/controller"
	"microdata/kemendagri/bumd/models"
	"microdata/kemendagri/bumd/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserHandler struct {
	Controller *controller.UserController
	Validate   *validator.Validate
}

func NewUserHandler(
	r fiber.Router,
	validator *validator.Validate,
	controller *controller.UserController,
) {
	handler := &UserHandler{
		Controller: controller,
		Validate:   validator,
	}

	// strict route
	rStrict := r.Group("user")
	rStrict.Get("/", handler.Index)
	rStrict.Get("/logout", handler.Logout)
	rStrict.Get("/profile", handler.Profile)
	rStrict.Get("/:id", handler.View)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
}

// Index func for get all user.
//
//	@Summary		get all user
//	@Description	get all user.
//	@ID				user-index
//	@Tags			User
//	@Produce		json
//	@Param			id_role	query		int					true	"Id Role"
//	@Param			page	query		int					false	"Halaman yang ditampilkan"
//	@Param			limit	query		int					false	"Jumlah data per halaman, maksimal 5 data per halaman"
//	@success		200		{object}	[]models.User		"Success"
//	@Failure		400		{object}	utils.RequestError	"Bad request"
//	@Failure		404		{object}	utils.RequestError	"Data not found"
//	@Failure		422		{array}		utils.RequestError	"Data validation failed"
//	@Failure		500		{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user [get]
func (h *UserHandler) Index(c *fiber.Ctx) error {
	idRole := c.QueryInt("id_role", 0)
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
		idRole,
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

// View func for get user by id.
//
//	@Summary		get user by id
//	@Description	get user by id.
//	@ID				user-view
//	@Tags			User
//	@Produce		json
//	@Param			id	path		int32				true	"Id untuk get data user"
//	@success		200	{object}	models.UserModel	"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user/{id} [get]
func (h *UserHandler) View(c *fiber.Ctx) error {
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

// Logout User func for logout.
//
//	@Summary		logout
//	@Description	user logout.
//	@ID				user-logout
//	@Tags			User
//	@Produce		json
//	@success		200	{object}	bool				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		403	{object}	utils.RequestError	"Forbidden"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user/logout [get]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	err := h.Controller.Logout(
		fmt.Sprintf("%v", claims["id_user"]),
	)
	if err != nil {
		return err
	}

	return c.JSON(true)
}

// Profile func for get profile info.
//
//	@Summary		user get profile info
//	@Description	get profile info.
//	@ID				user-profile
//	@Tags			User
//	@Produce		json
//	@success		200	{object}	models.UserDetail	"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user/profile [get]
func (h *UserHandler) Profile(c *fiber.Ctx) error {

	userModel, err := h.Controller.Profile(c.Locals("jwt").(*jwt.Token))
	if err != nil {
		return err
	}

	return c.JSON(userModel)
}

// Create func for create user.
//
//	@Summary		create user
//	@Description	create user.
//	@ID				user-create
//	@Tags			User
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			username	formData	string				true	"Username"
//	@Param			password	formData	string				true	"Password"
//	@Param			id_daerah	formData	int32				true	"Id Daerah"
//	@Param			id_role		formData	int32				true	"Id Role"
//	@Param			nama		formData	string				true	"Nama"
//	@Param			logo		formData	file				false	"Logo"
//	@Param			id_bumd		formData	int32				false	"Id BUMD"
//	@success		200			{object}	bool				"Success"
//	@Failure		400			{object}	utils.RequestError	"Bad request"
//	@Failure		403			{object}	utils.RequestError	"Forbidden"
//	@Failure		404			{object}	utils.RequestError	"Data not found"
//	@Failure		422			{array}		utils.RequestError	"Data validation failed"
//	@Failure		500			{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	payload := new(models.UserForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.Logo != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.Logo, maxFileSize, []string{"image/jpeg", "image/png", "image/jpg"}); err != nil {
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

// Update func for update user.
//
//	@Summary		update user
//	@Description	update user.
//	@ID				user-update
//	@Tags			User
//	@Produce		json
//	@Param			id			path		int32				true	"Id untuk update data user"
//	@Param			username	formData	string				true	"Username"
//	@Param			password	formData	string				true	"Password"
//	@Param			id_daerah	formData	int32				true	"Id Daerah"
//	@Param			id_role		formData	int32				true	"Id Role"
//	@Param			nama		formData	string				true	"Nama"
//	@Param			logo		formData	file				false	"Logo"
//	@Param			id_bumd		formData	int32				false	"Id BUMD"
//	@success		200			{object}	bool				"Success"
//	@Failure		400			{object}	utils.RequestError	"Bad request"
//	@Failure		403			{object}	utils.RequestError	"Forbidden"
//	@Failure		404			{object}	utils.RequestError	"Data not found"
//	@Failure		422			{array}		utils.RequestError	"Data validation failed"
//	@Failure		500			{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	payload := new(models.UserForm)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	if err := h.Validate.Struct(payload); err != nil {
		return err
	}

	if payload.Logo != nil {
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if err := utils.ValidateFile(payload.Logo, maxFileSize, []string{"image/jpeg", "image/png", "image/jpg"}); err != nil {
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

// Delete func for delete user.
//
//	@Summary		delete user
//	@Description	delete user.
//	@ID				user-delete
//	@Tags			User
//	@Produce		json
//	@Param			id	path		int32				true	"Id untuk delete data user"
//	@success		200	{object}	bool				"Success"
//	@Failure		400	{object}	utils.RequestError	"Bad request"
//	@Failure		403	{object}	utils.RequestError	"Forbidden"
//	@Failure		404	{object}	utils.RequestError	"Data not found"
//	@Failure		422	{array}		utils.RequestError	"Data validation failed"
//	@Failure		500	{object}	utils.RequestError	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/strict/user/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
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
