package dto

type QueryParam struct {
	SortBy   string `form:"sort" validate:"omitempty,oneof=server_name status ipv4 created_time"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
	Page     int    `form:"page" validate:"omitempty,gte=1"`
	PageSize int    `form:"page_size" validate:"omitempty,gte=1"`

	Name   string `form:"name"`
	Status string `form:"status" validate:"omitempty,oneof=on off"`
	IPv4   string `form:"ipv4"`
}
