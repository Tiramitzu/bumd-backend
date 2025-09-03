package models

type ResponseLogin struct {
	// Jwt token
	Token string `json:"token" xml:"token" example:"sdfsfsfsdfsfsdfsfsdfsf"`
	// Jwt refresh token
	RefreshToken string `json:"refresh_token" xml:"refresh_token" example:"sdfsfsfsdfsfsdfsfsdfsf"`
}

type LoginForm struct {
	Username string `json:"username"            form:"username" xml:"username" validate:"required" example:"warid"`
	Password string `json:"password"            form:"password" xml:"password" validate:"required" example:"123456"` // User password
}

type RefreshTokenForm struct {
	//JWT expired token
	Token string `json:"token" xml:"token" example:"xxxxx" validate:"required"`
}
