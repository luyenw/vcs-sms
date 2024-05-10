package service

import (
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
func (service *UserService) FindUserByID(id int) *entity.User {
	user := &entity.User{}
	err := service.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

func (service *UserService) CreateNewUser(username string, password string, scopes []entity.Scope) error {
	existScope := &entity.Scope{}
	scopeList := []entity.Scope{}
	for _, scope := range scopes {
		service.db.Where("name = ?", scope.Name).First(&existScope)
		if existScope.ID != 0 {
			scopeList = append(scopeList, *existScope)
		}
	}
	err := service.db.Create(&entity.User{Username: username, Password: password, Scopes: scopeList}).Error
	return err
}
func (service *UserService) FindByUsername(username string) entity.User {
	user := entity.User{}
	service.db.Where("username = ?", username).Preload("Scopes").First(&user)
	return user
}

func (service *UserService) UpdateUserScope(user *entity.User, scopes []dto.ScopeDTO) error {
	existScope := &entity.Scope{}
	scopeList := []entity.Scope{}
	for _, scope := range scopes {
		service.db.Where("name = ?", scope.Name).First(&existScope)
		if existScope.ID != 0 {
			scopeList = append(scopeList, *existScope)
		}
	}
	err := service.db.Model(&user).Association("Scopes").Replace(scopeList)
	return err
}
