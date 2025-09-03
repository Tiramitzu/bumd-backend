package utils

type RequestError struct {
	Code    int                   `json:"code" xml:"code" example:"422"`
	Message string                `json:"message" xml:"message" example:"Invalid email address"`
	Fields  []DataValidationError `json:"fields" xml:"fields"`
}

func (re RequestError) Error() string {
	return re.Message
}

type DataValidationError struct {
	Field   string `json:"field" xml:"field" example:"email"`
	Message string `json:"message" xml:"message" example:"Invalid email address"`
}

type GlobalError struct {
	Message string `json:"message" xml:"message" example:"invalid name"`
}

type LoginError struct {
	// Attempt   int    `json:"attempt" xml:"attempt" example:"3"`                            // sisa kesempatan login sebelum diblokir 5 menit
	// NextLogin int    `json:"next_login" xml:"next_login" example:"123233213"`              // unix timestamp UTC blokir login dibuka kembali
	Message string `json:"message" xml:"message" example:"invalid username or password"` // keterangan error
}

func (atp LoginError) Error() string {
	return atp.Message
}
