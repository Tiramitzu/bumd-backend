package models

type DashboardModel struct {
	TotalBumd  int64 `json:"total_bumd" xml:"total_bumd" example:"1"`
	TotalModal int64 `json:"total_modal" xml:"total_modal" example:"1"`
	TotalLaba  int64 `json:"total_laba" xml:"total_laba" example:"1"`
}
