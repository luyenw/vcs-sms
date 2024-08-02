package service

import (
	"vcs-sms/model/entity"
	"vcs-sms/repo"
)

type RegisteredMailService struct {
	DB repo.IDatabase
}

func NewRegisteredMailService(db repo.IDatabase) *RegisteredMailService {
	return &RegisteredMailService{
		DB: db,
	}
}

func (service *RegisteredMailService) GetAllRegisteredMails() []entity.RegisteredEmail {
	mails := []entity.RegisteredEmail{}
	service.DB.Find(&mails)
	return mails
}
