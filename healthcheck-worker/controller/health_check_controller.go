package controller

import (
	"encoding/json"
	"fmt"
	"healthcheck-worker/config/logger"
	"healthcheck-worker/config/mq"
	"healthcheck-worker/model/dto"
	"healthcheck-worker/model/entity"
	"healthcheck-worker/service"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type HealthCheckController struct {
	hcService     *service.HealthCheckService
	esService     *service.ESService
	serverService *service.ServerService
	cacheService  *service.CacheService
}

func NewHealthCheckController(hcService *service.HealthCheckService, esService *service.ESService, serverService *service.ServerService, cacheService *service.CacheService) *HealthCheckController {
	return &HealthCheckController{
		hcService:     hcService,
		esService:     esService,
		serverService: serverService,
		cacheService:  cacheService,
	}
}

func (h *HealthCheckController) allServerOffByDefault() {
	log := logger.NewLogger()
	var controllerOnStart sync.Once
	controllerOnStart.Do(func() {
		allServers := h.serverService.GetAllServers()
		allServersString, err := json.Marshal(allServers)
		if err != nil {
			log.Error(fmt.Sprintf("Error marshalling all servers: %v", err))
		}
		err = h.cacheService.Set("server:all", string(allServersString))
		if err != nil {
			log.Error(fmt.Sprintf("Error setting server:all in cache: %v", err))
			// log.Println("Error setting server:all in cache, ", err)
		}
		for _, server := range allServers {
			server.Status = 0
			server.LastUpdated = time.Now()
			if err := h.serverService.DB.Save(server).Error; err != nil {
				log.Error(fmt.Sprintf("Error saving server: %v", err))
				// log.Println(err)
			}
			doc := entity.ServerDoc{
				Server:    server,
				Timestamp: time.Now().UnixMilli(),
			}
			h.esService.InsertInBatch(doc)
		}
	})
}

func (h *HealthCheckController) HealthCheck() {
	log := logger.NewLogger()
	mailRequestChan := make(chan *kafka.Message, 100)
	c := mq.GetConsumer()

	go func() {
		for {
			msg, err := c.ReadMessage(-1)
			if err == nil {
				mailRequestChan <- msg
			}
		}
	}()
	go func() {
		workers := 10
		for i := 0; i < workers; i++ {
			go func() {
				for {
					select {
					case msg := <-mailRequestChan:
						healthCheckRequest := &dto.HealthcheckRequest{}
						if err := json.Unmarshal(msg.Value, healthCheckRequest); err != nil {
							log.Error(fmt.Sprintf("Error unmarshalling message: %v", err))
							continue
						}
						cachedServer := h.serverService.FindServerByIP(healthCheckRequest.Payload.AgentIP)

						if cachedServer == nil {
							cachedServer = &entity.Server{
								Name:        healthCheckRequest.Payload.AgentIP,
								IPv4:        healthCheckRequest.Payload.AgentIP,
								Status:      1,
								LastUpdated: time.Now(),
							}
							err := h.serverService.CreateServer(cachedServer)
							if err != nil {
								log.Error(fmt.Sprintf("Error creating server: %v", err))
								continue
							}
						}
						doc := entity.ServerDoc{
							Server:    *cachedServer,
							Timestamp: healthCheckRequest.Payload.Timestamp,
							Duration:  healthCheckRequest.Payload.Duration,
						}
						h.esService.InsertInBatch(doc)
					}
				}
			}()
		}
	}()
}
