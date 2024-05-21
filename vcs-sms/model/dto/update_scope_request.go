package dto

type ScopeDTO struct {
	Name string `json:"name"`
}

type UpdatePermissionRequest struct {
	Scopes []ScopeDTO `json:"scopes"`
}
