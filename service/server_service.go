package service

import (
	"log"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"

	"gorm.io/gorm"
)

type ServerService struct {
	DB *gorm.DB
}

func (service ServerService) GetAllServers() []entity.Server {
	servers := []entity.Server{}
	if err := service.DB.Find(&servers).Error; err != nil {
		log.Println(err)
		return []entity.Server{}
	}
	return servers
}

func NewServerService(db *gorm.DB) *ServerService {
	return &ServerService{
		DB: db,
	}
}

func (s *ServerService) CreateServer(server *entity.Server) error {
	err := s.DB.Create(server).Error
	return err
}

func (service *ServerService) GetServer(queryParam *dto.QueryParam) []entity.Server {
	servers := []entity.Server{}
	if queryParam.SortBy == "" {
		queryParam.SortBy = "created_time"
	}
	if queryParam.Order == "" {
		queryParam.Order = "desc"
	}
	if queryParam.PageSize == 0 {
		queryParam.PageSize = 10
	}
	if queryParam.Page == 0 {
		queryParam.Page = 1
	}
	err := service.DB.Order(queryParam.SortBy + " " + queryParam.Order).Limit(queryParam.PageSize).Offset((queryParam.Page - 1) * queryParam.PageSize).Find(&servers).Error
	if err != nil {
		return []entity.Server{}
	}
	return servers
}

func (s *ServerService) UpdateServer(server *entity.Server) error {
	err := s.DB.Save(server).Error
	return err
}

func (s *ServerService) DeleteServerById(id int) error {
	err := s.DB.Delete(&entity.Server{}, id).Error
	return err
}

func (s *ServerService) FindServerById(id int) *entity.Server {
	server := &entity.Server{}
	err := s.DB.First(server, id).Error
	if err != nil {
		return nil
	}
	return server
}
