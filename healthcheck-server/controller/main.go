package controller

import (
	"encoding/json"
	"fmt"
	"healthcheck-server/config/mq"
	"healthcheck-server/config/sql"
	"healthcheck-server/model"
	"healthcheck-server/service"
	"healthcheck-server/utils"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Payload struct {
	StatusCode        int    `json:"status_code"`
	ReceivedTimestamp string `json:"received_timestamp"`
	ResponseTimestamp string `json:"response_timestamp"`
	Message           string `json:"message"`
}

var statusMapping *utils.ConcurrentMap

func GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func Healthcheck(c *gin.Context) {
	var req model.HealthcheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := mq.GetProducer()
	topic := "healthcheck-topic"
	req.Payload.Timestamp = time.Now().UnixMilli()
	bytes, _ := json.Marshal(req)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	receivedTimestamp := time.Now().Format(time.RFC3339)
	responseTimestamp := time.Now().Add(1 * time.Second).Format(time.RFC3339)

	res := model.HealthcheckResponse{
		Payload: Payload{
			StatusCode:        200,
			ReceivedTimestamp: receivedTimestamp,
			ResponseTimestamp: responseTimestamp,
			Message:           "Healthcheck received successfully",
		},
	}

	statusMapping.Set(req.Payload.AgentIP, 1)

	c.JSON(http.StatusOK, res)
}

func UpdateServersStatus() {
	serverService := service.NewServerService(sql.GetPostgres())
	statusMapping = utils.NewConcurrentMap()
	for {
		ticker := time.NewTicker(1 * time.Minute)
		quit := make(chan struct{})
		select {
		case <-ticker.C:
			fmt.Println("Updating servers status")
			serverService.UpdateServersOn(statusMapping.Items())
		case <-quit:
			ticker.Stop()
		}
	}
}
