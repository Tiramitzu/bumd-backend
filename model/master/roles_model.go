package models

type RolesModel struct {
	ID   int    `json:"id" xml:"id"`
	Nama string `json:"nama" xml:"nama"`
}

type RolesForm struct {
	Nama string `json:"nama" xml:"nama"`
}
