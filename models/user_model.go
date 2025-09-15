package models

import (
	"encoding/json"
	"mime/multipart"

	"github.com/google/uuid"
)

type User struct {
	KodeDDN         string    `json:"kode_ddn"          xml:"kode_ddn"`
	KodeProvinsi    string    `json:"kode_provinsi"    xml:"kode_provinsi"`
	SubDomainDaerah string    `json:"sub_domain_daerah" xml:"sub_domain_daerah"`
	IdUser          int64     `json:"id_user"           xml:"id_user"`
	IdDaerah        int32     `json:"id_daerah"         xml:"id_daerah"`
	IdRole          int32     `json:"id_role"           xml:"id_role"`
	IdBumd          uuid.UUID `json:"id_bumd"           xml:"id_bumd"`
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
	Username   string    `json:"username"          xml:"username"`
	NamaUser   string    `json:"nama_user"         xml:"nama_user"`
	NamaRole   string    `json:"nama_role"         xml:"nama_role"`
	NamaBumd   string    `json:"nama_bumd"         xml:"nama_bumd"`
	NamaDaerah string    `json:"nama_daerah"       xml:"nama_daerah"`
	IdUser     int64     `json:"id_user"           xml:"id_user"`
	IdDaerah   int32     `json:"id_daerah"         xml:"id_daerah"`
	IdRole     int32     `json:"id_role"           xml:"id_role"`
	IdBumd     uuid.UUID `json:"id_bumd"           xml:"id_bumd"`
	KodeDDN    string    `json:"kode_ddn"          xml:"kode_ddn"`
	KodeProp   string    `json:"kode_prop"    xml:"kode_prop"`
	SubDomain  string    `json:"sub_domain" xml:"sub_domain"`
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
	Username string `json:"username"          xml:"username"`
	Nama     string `json:"nama"              xml:"nama"`
	Logo     string `json:"logo"              xml:"logo"`
	Role     string `json:"role"              xml:"role"`
	IdUser   int64  `json:"id_user"           xml:"id_user"`
	IdDaerah int32  `json:"id_daerah"         xml:"id_daerah"`
	IdRole   int32  `json:"id_role"           xml:"id_role"`
	IdBumd   int32  `json:"id_bumd"           xml:"id_bumd"`
}

type UserForm struct {
	Username string                `json:"username" xml:"username" form:"username"`
	Password string                `json:"password" xml:"password" form:"password"`
	Nama     string                `json:"nama" xml:"nama" form:"nama"`
	Logo     *multipart.FileHeader `json:"logo" xml:"logo" form:"logo"`
	IdDaerah int32                 `json:"id_daerah" xml:"id_daerah" form:"id_daerah"`
	IdRole   int32                 `json:"id_role" xml:"id_role" form:"id_role"`
	IdBumd   int32                 `json:"id_bumd" xml:"id_bumd" form:"id_bumd"`
}
