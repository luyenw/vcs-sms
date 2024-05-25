package service

import (
	"context"
	"encoding/json"
	"fmt"
	"vcs-sms/model/dto"
	"vcs-sms/repo"

	"github.com/elastic/go-elasticsearch/v8/esutil"

	"vcs-sms/config/elasticsearch"
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
	res, err := client.UptimeCheck(context.Background(), &pb.UptimeCheckRequest{
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
	bi    esutil.BulkIndexer
}

func NewESService(escli repo.ElasticRepo) *ESService {
	log := logger.NewLogger()
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: elasticsearch.GetESClient(),
		Index:  "vcs-sms",
	})
	if err != nil {
		log.Error(fmt.Sprintf("Error creating new bulk indexer: %s", err))
		return nil
	}
	return &ESService{
		escli: escli,
		bi:    bi,
	}
}

func (service *ESService) CalculateUptime(startMils int64, endMils int64) []dto.ServerUptime {
	srv := NewRpcService()
	uptime := srv.UptimeCheck(startMils, endMils)
	return uptime
}
