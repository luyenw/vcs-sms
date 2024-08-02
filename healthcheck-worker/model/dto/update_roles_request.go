package dto

type UpdateRolesRequest struct {
	Roles []string `json:"roles" binding:"required"`
}
