package models

type SysConfigModel struct {
	ID                      int32 `json:"id" xml:"id"`
	OtpExpiredMinutes       int32 `json:"otp_expired_minutes" xml:"otp_expired_minutes"`
	JwtExpiredMinutes       int32 `json:"jwt_expired_minutes" xml:"jwt_expired_minutes"`
	RefreshTokenExpiredHour int32 `json:"refresh_token_expired_hour" xml:"refresh_token_expired_hour"`
}

type SysConfigForm struct {
	OtpExpiredMinutes       int32 `json:"otp_expired_minutes" xml:"otp_expired_minutes"`
	JwtExpiredMinutes       int32 `json:"jwt_expired_minutes" xml:"jwt_expired_minutes"`
	RefreshTokenExpiredHour int32 `json:"refresh_token_expired_hour" xml:"refresh_token_expired_hour"`
}
