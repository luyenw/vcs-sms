package main

import (
	"encoding/json"
	"fmt"
	"mail-service/config"
	"mail-service/config/logger"
	"mail-service/config/mq"
	"mail-service/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type mailRequest struct {
	To   []string
	Body string
}

func main() {
	config.InitConfig()
	log := logger.NewLogger()
	log.Info("Server started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Fatal("Received signal, shutting down", zap.String("signal", sig.String()))
		os.Exit(0)
	}()

	c := mq.GetConsumer()
	fmt.Println(c)
	mailService := service.NewMailService()

	mailRequestChan := make(chan *kafka.Message, 100)
	go func() {
		workers := 10
		for i := 0; i < workers; i++ {
			go func() {
				for {
					select {
					case msg := <-mailRequestChan:
						mailRequest := &mailRequest{}
						if err := json.Unmarshal(msg.Value, mailRequest); err != nil {
							log.Error(fmt.Sprintf("Error unmarshalling message: %v", err))
							continue
						}
						if err := mailService.SendEmail(mailRequest.To, mailRequest.Body); err != nil {
							log.Error(fmt.Sprintf("Error sending email to %v: %v", mailRequest.To, err))
							continue
						}
						log.Info(fmt.Sprintf("Sent email to %v", mailRequest.To))
					}
				}
			}()
		}
	}()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			mailRequestChan <- msg
		}
	}
}
