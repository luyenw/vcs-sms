package dto

type QueryParam struct {
	Filters  string `form:"filters"`
	SortBy   string `form:"sort"`
	Order    string `form:"order"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
