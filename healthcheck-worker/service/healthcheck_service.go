package service

import (
	"fmt"
	"healthcheck-worker/config/logger"
	"time"

	"github.com/go-ping/ping"
)

type HealthCheckService struct {
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (h *HealthCheckService) HealthCheck(addr string) int {
	select {
	case <-time.After(10 * time.Second):
		return 0
	case status := <-icmpEcho(addr):
		return status
	}
}

func icmpEcho(ipv4 string) chan int {
	log := logger.NewLogger()
	result := make(chan int)
	go func() {
		pinger, err := ping.NewPinger(ipv4)
		// pinger.SetPrivileged(true)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating pinger: %v", err))
			result <- 0
		}
		pinger.Count = 3
		pinger.OnFinish = func(s *ping.Statistics) {
			if s.PacketsRecv != 0 {
				result <- 1
			}
		}
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			log.Error(fmt.Sprintf("Error running pinger: %v", err))
			result <- 0

		}
	}()
	return result
}
