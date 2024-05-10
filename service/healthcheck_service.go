package service

import (
	"log"
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
	case <-time.After(5 * time.Second):
		return 0
	case <-icmpEcho(addr):
		return 1
	}
}

func icmpEcho(ipv4 string) chan int {
	result := make(chan int)
	go func() {
		pinger, err := ping.NewPinger(ipv4)
		pinger.SetPrivileged(true)
		if err != nil {
			log.Println(err.Error())
		}
		pinger.Count = 3
		pinger.OnFinish = func(s *ping.Statistics) {
			if s.PacketsRecv != 0 {
				result <- 1
			}
		}
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return result
}