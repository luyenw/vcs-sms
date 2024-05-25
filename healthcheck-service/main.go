package main

import (
	"healthcheck-service/config/cache"
	"healthcheck-service/config/elasticsearch"
	"healthcheck-service/config/logger"
	"healthcheck-service/config/rpc"
	"healthcheck-service/config/sql"
	"healthcheck-service/controller"
	"healthcheck-service/service"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vcs-sms/repo"

	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()
	log.Info("Server started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Fatal("Received signal, shutting down", zap.String("signal", sig.String()))
		os.Exit(0)
	}()

	healthCheckController := controller.NewHealthCheckController(
		service.NewHealthCheckService(),
		service.NewESService(&repo.ESClient{Client: elasticsearch.GetESClient()}),
		service.NewServerService(sql.GetPostgres()),
		service.NewCacheService(cache.GetRedis()),
	)

	go func() {
		rpc.ConfigGRPC()
	}()

	go func() {
		log := logger.NewLogger()
		for {
			select {
			case <-time.After(1 * time.Second):
				log.Info("Certainly! Here's the text with `\n` added to represent line breaks:				In the vast expanse of the cosmos, where stars twinkle like distant memories and galaxies swirl in a cosmic dance, there exists a realm of infinite possibilities. \nIt is a realm where time stretches and contracts like a rubber band, where the laws of physics bend and twist like a river flowing through the fabric of space-time.\n\nIn this cosmic tapestry, planets are born from the remnants of dying stars, and life emerges from the primordial soup of ancient oceans. \nIt is a place where the echoes of ancient civilizations reverberate through the void, whispering secrets of forgotten worlds.\n\nBut amidst the beauty and wonder of the cosmos, there is also darkness. Black holes lurk in the depths of space, swallowing everything that dares to venture too close. \nNebulas glow with an eerie light, casting shadows that dance across the surface of distant moons.\n\nYet even in the darkest corners of the universe, there is hope. For where there is darkness, there is also light. \nAnd where there is light, there is life.\n\nAcross the vast expanse of space, civilizations rise and fall like waves crashing upon the shore. \nEmpires expand and contract, wars rage and peace reigns, and through it all, the cosmos continues its eternal dance.\n\nAnd so, as we gaze up at the stars on a clear night, let us remember that we are but a small part of something much greater. \nWe are children of the cosmos, born from stardust and destined to return to the stars. \nAnd in that realization, we find both humility and awe, for we are but fleeting travelers in the grand adventure of the universe.")
			}
		}
	}()
	done := make(chan interface{})
	healthCheckController.HealthCheck()
	<-done
}
