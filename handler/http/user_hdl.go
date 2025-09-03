package http

import (
	"fmt"
	"microdata/kemendagri/bumd/controller"
	models "microdata/kemendagri/bumd/model"
	"microdata/kemendagri/bumd/utils"
	"strings"

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
	rStrict.Get("/logout", handler.Logout)
	rStrict.Get("/profile", handler.Profile)
	rStrict.Post("/", handler.Create)
	rStrict.Put("/:id", handler.Update)
	rStrict.Delete("/:id", handler.Delete)
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
//	@Param			id_daerah	formData	int					true	"Id Daerah"
//	@Param			id_role		formData	int					true	"Id Role"
//	@Param			nama		formData	string				true	"Nama"
//	@Param			logo		formData	string				false	"Logo"
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
		if !strings.Contains(payload.Logo.Header.Get("Content-Type"), "image") {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: "Logo harus berupa gambar",
			}
		}
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if payload.Logo.Size > maxFileSize {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: "Ukuran logo maksimal 2 MB",
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
//	@Param			id			path		int					true	"Id untuk update data user"
//	@Param			username	formData	string				true	"Username"
//	@Param			password	formData	string				true	"Password"
//	@Param			id_daerah	formData	int					true	"Id Daerah"
//	@Param			id_role		formData	int					true	"Id Role"
//	@Param			nama		formData	string				true	"Nama"
//	@Param			logo		formData	string				false	"Logo"
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
		if !strings.Contains(payload.Logo.Header.Get("Content-Type"), "image") {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: "Logo harus berupa gambar",
			}
		}
		const maxFileSize = 2 * 1024 * 1024 // 2 MB
		if payload.Logo.Size > maxFileSize {
			return utils.RequestError{
				Code:    fiber.StatusBadRequest,
				Message: "Ukuran logo maksimal 2 MB",
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
//	@Param			id	path		int					true	"Id untuk delete data user"
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
