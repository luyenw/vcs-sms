package main

import (
	"encoding/json"
	"fmt"
	"mail-service/config/logger"
	"mail-service/config/mq"
	"mail-service/service"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	c := mq.GetConsumer()
	mailService := service.NewMailService()
	type mailRequest struct {
		To   []string
		Body string
	}
	for {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			mailRequest := &mailRequest{}
			if e := json.Unmarshal(msg.Value, mailRequest); err != nil {
				fmt.Println(e)
				log.Error(fmt.Sprintf("Error unmarshalling message: %v", e))
				continue
			}
			if e := mailService.SendEmail(mailRequest.To, mailRequest.Body); e != nil {
				fmt.Println(e)
				log.Error(fmt.Sprintf("Error sending email to %v: %v", mailRequest.To, e))
				continue
			}
			log.Info(fmt.Sprintf("Sent email to %v", mailRequest.To))
		} else {
			// fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			// log.Error(fmt.Sprintf("Consumer error: %v (%v)\n", err, msg))
		}
	}
}
