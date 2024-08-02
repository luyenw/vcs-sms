package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"healthcheck-server/config/elasticsearch"
	"healthcheck-server/model/dto"
	"healthcheck-server/repo"

	"healthcheck-server/config/logger"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
)

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
		log.Error(fmt.Sprintf("Error creating the indexer: %s", err))
		log.Fatal("Shutting down")
		// log.Fatalf("Error creating the indexer: %s", err)
	}
	return &ESService{
		escli: escli,
		bi:    bi,
	}
}

func (service ESService) InsertInBatch(doc interface{}) {
	log := logger.NewLogger()
	data, err := json.Marshal(doc)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling the document: %s", err))
		return
	}
	err = service.bi.Add(context.Background(), esutil.BulkIndexerItem{
		Action:     "create",
		DocumentID: uuid.New().String(),
		Body:       bytes.NewReader(data),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			log.Info(fmt.Sprintf("Document added to the indexer: %s", res.Result))
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			// log.Println(err)
		},
	})

	if err != nil {
		// log.Println(err)
		log.Error(fmt.Sprintf("Error adding the document to the indexer: %s", err))
	}
}

func (service *ESService) CalculateUptime(startMils int64, endMils int64) []dto.ServerUptime {
	log := logger.NewLogger()
	if startMils < 0 || endMils < 0 {
		return []dto.ServerUptime{}
	}
	var query = fmt.Sprintf(`
	{
	"size": 0,
	"query": {"range": {"timestamp": {"gte": %d, "lte": %d}}},
      "aggs": {
        "by_server": {
          "terms": {
            "field": "server.ID", 
            "size": 10000 
          },
          "aggs": {
            "total_duration": {
              "sum": {
                "field": "duration"
              }
            },
            "min_timestamp": {
              "min": {
                "field": "timestamp"
              }
            },
            "max_timestamp": {
              "max": {
                "field": "timestamp"
              }
            },
            "uptime_avg": {
              "bucket_script": {
                "buckets_path": {
                  "totalDuration": "total_duration",
                  "minTimestamp": "min_timestamp",
                  "maxTimestamp": "max_timestamp"
                },
			"script": "params.minTimestamp == null ? 1 : Math.min(1, params.totalDuration / (params.maxTimestamp / 1000) - (params.minTimestamp / 1000)))"              }
            }
          }
        }
      }
    }
    `, startMils, endMils)
	res, err := service.escli.Query(query)
	var response dto.Response
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		log.Error(fmt.Sprintf("Error decoding the response: %s", err))
		return []dto.ServerUptime{}
	}

	return response.Aggregtions.Server.Buckets
}
