package dto

import "vcs-sms/model/entity"

type UserDTO struct {
	Username string      `json:"Username"`
	Role     entity.Role `json:"Role"`
}

type UserEntity struct {
	*entity.User
}

func (u UserEntity) ToDTO() UserDTO {
	return UserDTO{
		Username: u.Username,
		Role:     u.Role,
	}
}
