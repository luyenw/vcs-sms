package dto

import "vcs-sms/model/entity"

type UpdateRoleRequest struct {
	Role entity.Role `json:"role" binding:"required"`
}
