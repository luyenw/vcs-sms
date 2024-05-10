package controller

import (
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
}

func NewHealthCheckController(hcService *service.HealthCheckService, esService *service.ESService, serverService *service.ServerService) *HealthCheckController {
	return &HealthCheckController{
		hcService:     hcService,
		esService:     esService,
		serverService: serverService,
	}
}

func (h *HealthCheckController) allServerOffByDefault(servers []entity.Server) {
	var controllerOnStart sync.Once
	controllerOnStart.Do(func() {
		for _, server := range servers {
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

	h.allServerOffByDefault(h.serverService.GetAllServers())

	jobs := make(chan int, 100)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			servers := h.serverService.GetAllServers()
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
					server := h.serverService.FindServerById(job)
					server.Status = h.hcService.HealthCheck(server.IPv4)
					server.LastUpdated = time.Now()
					if err := h.serverService.DB.Save(server).Error; err != nil {
						log.Println(err)
					}
					doc := entity.ServerDoc{
						Server:    *server,
						Timestamp: time.Now().UnixMilli(),
					}
					h.esService.InsertInBatch(doc)
				}
			}()

		}
	}()
}
