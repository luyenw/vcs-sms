package dto

type QueryParam struct {
	Filters  string `form:"filters" validate:"string"`
	SortBy   string `form:"sort" validate:"string"`
	Order    string `form:"order" validate:"oneof=asc desc"`
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1"`
}
