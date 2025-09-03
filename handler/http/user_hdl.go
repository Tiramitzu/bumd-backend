package http

import (
	"fmt"
	"microdata/kemendagri/bumd/controller"

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
