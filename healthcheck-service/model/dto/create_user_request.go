package dto

import "healthcheck-service/model/entity"

type CreateUserRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Scopes   []entity.Scope `json:"scopes" binding:"required"`
}
