package dto

import "vcs-sms/model/entity"

type CreateUserRequest struct {
	Username string      `json:"username" binding:"required"`
	Password string      `json:"password" binding:"required"`
	Role     entity.Role `json:"role" binding:"required"`
}
