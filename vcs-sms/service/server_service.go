package service

import (
	"fmt"
	"vcs-sms/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"
	"vcs-sms/repo"
)

type IServerService interface {
	CreateServer(server *entity.Server) error
	DeleteServerById(id int) error
	FindServerById(id int) *entity.Server
	GetAllServers() []entity.Server
	GetServer(queryParam *dto.QueryParam) []entity.Server
	UpdateServer(server *entity.Server) error
}

type ServerService struct {
	DB repo.IDatabase
}

func NewServerService(db repo.IDatabase) *ServerService {
	return &ServerService{
		DB: db,
	}
}

func (service ServerService) GetAllServers() []entity.Server {
	log := logger.NewLogger()
	servers := []entity.Server{}
	if err := service.DB.Find(&servers).Error; err != nil {
		log.Error(fmt.Sprintf("Error getting all servers: %s", err))
		return []entity.Server{}
	}
	return servers
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
	if err := service.DB.Order(queryParam.SortBy + " " + queryParam.Order).Offset((queryParam.Page - 1) * queryParam.PageSize).Limit(queryParam.PageSize).Find(&servers).Error; err != nil {
		fmt.Println(err.Error())
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
