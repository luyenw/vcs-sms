package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"vcs-sms/model/entity"
	"vcs-sms/service"
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
	var controllerOnStart sync.Once
	controllerOnStart.Do(func() {
		allServers := h.serverService.GetAllServers()
		allServersString, err := json.Marshal(allServers)
		if err != nil {
			log.Println(err)
		}
		err = h.cacheService.Set("server:all", string(allServersString))
		if err != nil {
			log.Println("Error setting server:all in cache, ", err)
		}
		for _, server := range allServers {
			server.Status = 0
			server.LastUpdated = time.Now()
			if err := h.serverService.DB.Save(server).Error; err != nil {
				log.Println(err)
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
	h.allServerOffByDefault()

	jobs := make(chan int, 100)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			serversString, err := h.cacheService.Get("server:all")
			if err != nil {
				log.Println("Error getting server:all from cache, ", err)
			}
			servers := []entity.Server{}
			if err := json.Unmarshal([]byte(serversString), &servers); err != nil {
				servers = h.serverService.GetAllServers()
				serversString, _ := json.Marshal(servers)
				h.cacheService.Set("server:all", string(serversString))
			}
			select {
			case <-ticker.C:
				go func() {
					for _, server := range servers {
						jobs <- int(server.ID)
					}
				}()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		numWokers := 5
		for w := 0; w < numWokers; w++ {
			go func() {
				for job := range jobs {
					serverString, err := h.cacheService.Get(fmt.Sprintf("server:%d", job))
					if err != nil {
						log.Println("Error getting server:all from cache, ", err)
					}
					cachedServer := &entity.Server{}
					if err := json.Unmarshal([]byte(serverString), cachedServer); err != nil {
						log.Println(err)
						server := h.serverService.FindServerById(job)
						serverString, _ := json.Marshal(server)
						h.cacheService.Set(fmt.Sprintf("server:%d", job), string(serverString))
					}
					cachedServer.Status = h.hcService.HealthCheck(cachedServer.IPv4)
					cachedServer.LastUpdated = time.Now()
					if err := h.serverService.DB.Table("servers").Where("id=?", job).Update("status", cachedServer.Status).Update("last_updated", cachedServer.LastUpdated).Error; err != nil {
						log.Println(err)
					}
					doc := entity.ServerDoc{
						Server:    *cachedServer,
						Timestamp: time.Now().UnixMilli(),
					}
					h.esService.InsertInBatch(doc)
				}
			}()

		}
	}()
}
