package models

type SiteTestDbForm struct {
	DbConn      string `json:"db_conn" xml:"db_conn" validate:"required"`
	QueryString string `json:"query_string" xml:"query_string" validate:"required"`
	/*DbScheme string `json:"db_scheme" xml:"db_scheme" validate:"required"`
	DbTable  string `json:"db_table" xml:"db_table" validate:"required"`
	DbColumn string `json:"db_column" xml:"db_column" validate:"required"`*/
}
