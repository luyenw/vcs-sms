package service

import (
	"fmt"
	"vcs-sms/model/entity"
	"vcs-sms/repo"
)

type IUserService interface {
	CreateNewUser(username string, password string, Role entity.Role) error
	FindByUsername(username string) entity.User
	FindUserByID(id int) *entity.User
	UpdateUserRole(user *entity.User, role *entity.Role) error
	FindRoleByName(name string) (*entity.Role, error)
}

type UserService struct {
	db repo.IDatabase
}

func NewUserService(db repo.IDatabase) *UserService {
	return &UserService{db: db}
}
func (service *UserService) FindUserByID(id int) *entity.User {
	user := &entity.User{}
	err := service.db.Where("id = ?", id).Preload("Role.Scopes").First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

func (service *UserService) CreateNewUser(username string, password string, role entity.Role) error {
	r := entity.Role{}
	service.db.Where("name = ?", role.Name).Preload("Scopes").First(&r)
	fmt.Println(r)
	err := service.db.Create(&entity.User{Username: username, Password: password, Role: r}).Error
	return err
}
func (service *UserService) FindByUsername(username string) entity.User {
	user := entity.User{}
	if err := service.db.Where("username = ?", username).Preload("Role.Scopes").First(&user).Error; err != nil {
		return entity.User{}
	}
	return user
}

func (service *UserService) UpdateUserRole(user *entity.User, role *entity.Role) error {
	user.RoleID = role.ID
	err := service.db.Save(&user).Error
	return err
}

func (service *UserService) FindRoleByName(name string) (*entity.Role, error) {
	role := &entity.Role{}
	err := service.db.Where("name = ?", name).First(role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}
