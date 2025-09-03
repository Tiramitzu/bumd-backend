package http

import (
	"microdata/kemendagri/bumd/controller"
	"microdata/kemendagri/bumd/handler/http/http_util"
	models "microdata/kemendagri/bumd/model"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Controller *controller.AuthController
	Validate   *validator.Validate
}

func NewAuthHandler(app *fiber.App, controller *controller.AuthController, vld *validator.Validate) {
	handler := &AuthHandler{
		Controller: controller,
		Validate:   vld,
	}

	// public route
	rPub := app.Group("/auth")
	rPub.Post("/login", handler.Login)
}

// Login func for login.
//
//	@Summary		user login
//	@Description	Login to get JWT token and refresh token.
//	@ID				auth-login
//	@Tags			Auth
//	@Accept			json
//	@Param			payload	body	models.LoginForm	true	"Login payload"
//	@Produce		json
//	@Success		200	{object}	http_util.JSONResultLogin	"Login Success, jwt provided"
//	@Failure		400	{object}	utils.RequestError			"Bad request"
//	@Failure		403	{object}	utils.LoginError			"Login forbidden"
//	@Failure		404	{object}	utils.RequestError			"Data not found"
//	@Failure		422	{array}		utils.RequestError			"Data validation failed"
//	@Failure		500	{object}	utils.RequestError			"Server error"
//	@Router			/auth/login [post]
func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	formModel := new(models.LoginForm)
	if err := c.BodyParser(formModel); err != nil {
		return err
	}

	token, refreshToken, err := ah.Controller.Login(*formModel)
	if err != nil {
		return err
	}

	return c.JSON(
		http_util.JSONResultLogin{Token: token, RefreshToken: refreshToken},
	)
}
