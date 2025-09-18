package models

import (
	"encoding/json"
	"mime/multipart"

	"github.com/google/uuid"
)

type User struct {
	KodeDDN         string    `json:"kode_ddn"          xml:"kode_ddn" example:"123"`
	KodeProvinsi    string    `json:"kode_provinsi"     xml:"kode_provinsi" example:"11.01"`
	SubDomainDaerah string    `json:"sub_domain_daerah" xml:"sub_domain_daerah" example:"jakarta"`
	IdUser          int64     `json:"id_user"           xml:"id_user" example:"1"`
	IdDaerah        int32     `json:"id_daerah"         xml:"id_daerah" example:"1"`
	IdRole          int32     `json:"id_role"           xml:"id_role" example:"1"`
	IdBumd          uuid.UUID `json:"id_bumd"           xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
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
	Username   string    `json:"username"          xml:"username" example:"admin"`
	NamaUser   string    `json:"nama_user"         xml:"nama_user" example:"Admin"`
	NamaRole   string    `json:"nama_role"         xml:"nama_role" example:"Admin"`
	NamaBumd   string    `json:"nama_bumd"         xml:"nama_bumd" example:"BUMD"`
	NamaDaerah string    `json:"nama_daerah"       xml:"nama_daerah" example:"Daerah"`
	IdUser     int64     `json:"id_user"           xml:"id_user" example:"1"`
	IdDaerah   int32     `json:"id_daerah"         xml:"id_daerah" example:"1"`
	IdRole     int32     `json:"id_role"           xml:"id_role" example:"1"`
	IdBumd     uuid.UUID `json:"id_bumd"           xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
	KodeDDN    string    `json:"kode_ddn"          xml:"kode_ddn" example:"11.01"`
	KodeProp   string    `json:"kode_prop"    xml:"kode_prop" example:"11.01"`
	SubDomain  string    `json:"sub_domain" xml:"sub_domain" example:"jakarta"`
	Logo       string    `json:"logo" xml:"logo" example:"/path/to/file.png"`
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

type UserModel struct {
	Username string    `json:"username"          xml:"username" example:"admin"`
	Nama     string    `json:"nama"              xml:"nama" example:"Admin"`
	Logo     string    `json:"logo"              xml:"logo" example:"/path/to/file.png"`
	Role     string    `json:"role"              xml:"role" example:"Admin"`
	IdUser   int64     `json:"id_user"           xml:"id_user" example:"1"`
	IdDaerah int32     `json:"id_daerah"         xml:"id_daerah" example:"1"`
	IdRole   int32     `json:"id_role"           xml:"id_role" example:"1"`
	IdBumd   uuid.UUID `json:"id_bumd"           xml:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type UserForm struct {
	Username string                `json:"username" xml:"username" form:"username" validate:"required" example:"admin"`
	Password string                `json:"password" xml:"password" form:"password" validate:"required" example:"password"`
	Nama     string                `json:"nama" xml:"nama" form:"nama" validate:"required" example:"Admin"`
	Logo     *multipart.FileHeader `json:"logo" xml:"logo" form:"logo"`
	IdDaerah int32                 `json:"id_daerah" xml:"id_daerah" form:"id_daerah" validate:"required" example:"1"`
	IdRole   int32                 `json:"id_role" xml:"id_role" form:"id_role" validate:"required" example:"1"`
	IdBumd   uuid.UUID             `json:"id_bumd" xml:"id_bumd" form:"id_bumd" example:"123e4567-e89b-12d3-a456-426614174000" default:"00000000-0000-0000-0000-000000000000"`
}
