package service

import (
	"context"
	"encoding/json"
	"fmt"
	"vcs-sms/model/dto"
	"vcs-sms/repo"

	"vcs-sms/config/logger"
	"vcs-sms/config/rpc"

	pb "vcs-sms/grpc"
)

type IRpcService interface {
	UptimeCheck(startMilis int64, endMilis int64) []dto.ServerUptime
}

type RpcService struct{}

func NewRpcService() *RpcService {
	return &RpcService{}
}
func (s *RpcService) UptimeCheck(startMilis int64, endMilis int64) []dto.ServerUptime {
	log := logger.NewLogger()
	client := rpc.GetRpcClient()
	res, err := (*client).UptimeCheck(context.Background(), &pb.UptimeCheckRequest{
		StartTime: startMilis,
		EndTime:   endMilis,
	})
	if err != nil {
		log.Error(err.Error())
		return []dto.ServerUptime{}
	}
	var uptime []dto.ServerUptime
	json.Unmarshal([]byte(res.Response), &uptime)
	return uptime
}

type ESService struct {
	escli repo.ElasticRepo
}

func NewESService(escli repo.ElasticRepo) *ESService {
	return &ESService{
		escli: escli,
	}
}

func (service *ESService) CalculateUptime(startMils int64, endMils int64) []dto.ServerUptime {
	srv := NewRpcService()
	uptime := srv.UptimeCheck(startMils, endMils)
	fmt.Println("Start time: ", startMils, " - End time: ", endMils)
	fmt.Println("Uptime: ", uptime)
	return uptime
}
