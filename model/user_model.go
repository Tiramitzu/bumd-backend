package models

import (
	"encoding/json"
	"mime/multipart"
)

type User struct {
	IdUser   uint64 `json:"id_user"           xml:"id_user"`
	IdDaerah uint64 `json:"id_daerah"         xml:"id_daerah"`
	IdRole   int    `json:"id_role"           xml:"id_role"`
}

// FromJSON decode json to user struct
func (u *User) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, u)
}

// ToJSON encode user struct to json
func (u *User) ToJSON() []byte {
	str, _ := json.Marshal(u)
	return str
}

type UserDetail struct {
	Username string `json:"username"          xml:"username"`
	NamaUser string `json:"nama_user"         xml:"nama_user"`
	NamaRole string `json:"nama_role"         xml:"nama_role"`
}

// FromJSON decode json to user struct
func (u *UserDetail) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, u)
}

// ToJSON encode user struct to json
func (u *UserDetail) ToJSON() []byte {
	str, _ := json.Marshal(u)
	return str
}

type UserForm struct {
	Username string                `json:"username" xml:"username"`
	Password string                `json:"password" xml:"password"`
	IdDaerah int                   `json:"id_daerah" xml:"id_daerah"`
	IdRole   int                   `json:"id_role" xml:"id_role"`
	Nama     string                `json:"nama" xml:"nama"`
	Logo     *multipart.FileHeader `json:"logo" xml:"logo"`
}
