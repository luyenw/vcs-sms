package service

import (
	"vcs-sms/config/sql"
	"vcs-sms/model/entity"

	"gorm.io/gorm"
)

type RegisteredMailService struct {
	DB *gorm.DB
}

func NewRegisteredMailService() *RegisteredMailService {
	return &RegisteredMailService{
		DB: sql.GetPostgres(),
	}
}

func (service *RegisteredMailService) GetAllRegisteredMails() []entity.RegisteredEmail {
	mails := []entity.RegisteredEmail{}
	service.DB.Find(&mails)
	return mails
}
