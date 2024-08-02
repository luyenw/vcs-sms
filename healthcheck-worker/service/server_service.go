package service

import (
	"fmt"
	"healthcheck-worker/config/logger"
	"healthcheck-worker/model/entity"
	"healthcheck-worker/repo"
)

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
		log.Error(fmt.Sprintf("Error finding all servers: %v", err))
		return []entity.Server{}
	}
	return servers
}

func (s *ServerService) UpdateServer(server *entity.Server) error {
	log := logger.NewLogger()
	var err error
	if err = s.DB.Save(server).Error; err != nil {
		log.Error(fmt.Sprintf("Error saving server: %v", err))
	}
	return err
}

func (s *ServerService) FindServerById(id int) *entity.Server {
	log := logger.NewLogger()
	server := &entity.Server{}
	err := s.DB.First(server, id).Error
	if err != nil {
		log.Error(fmt.Sprintf("Error finding server by id: %v", err))
		return nil
	}
	return server
}

func (s *ServerService) FindServerByIP(ipv4 string) *entity.Server {
	log := logger.NewLogger()
	server := &entity.Server{}
	err := s.DB.First(server, "ipv4=?", ipv4).Error
	if err != nil {
		log.Error(fmt.Sprintf("Error finding server by ipv4: %v", err))
		return nil
	}
	return server
}

func (s *ServerService) CreateServer(server *entity.Server) error {
	log := logger.NewLogger()
	var err error
	if err = s.DB.Create(server).Error; err != nil {
		log.Error(fmt.Sprintf("Error creating server: %v", err))
	}
	return err
}
