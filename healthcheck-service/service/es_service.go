package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mail-service/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/repo"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"

	"vcs-sms/config/elasticsearch"
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
		"server": {
		"terms": {"field": "server.ID", "size": 100000},
		"aggs": {
		"uptime": {
		"scripted_metric": {
		"init_script": "state.docs = [];",
		"map_script": "def dto = [:]; dto.timestamp = doc['timestamp'].value; dto.status = doc['server.Status'].value; state.docs.add(dto);",
		"combine_script": "def sortedDocs = state.docs.stream().sorted(Comparator.comparingLong(doc -> doc.timestamp.toInstant().toEpochMilli())).collect(Collectors.toList()); def uptimes = []; for (int i = 0; i < sortedDocs.size() - 1; i++) { if(sortedDocs[i+1].status == 1){ def timestamp1 = sortedDocs[i].timestamp.toInstant().toEpochMilli(); def timestamp2 = sortedDocs[i + 1].timestamp.toInstant().toEpochMilli(); def up = timestamp2 - timestamp1; uptimes.add(up); } } uptimes.add(sortedDocs[sortedDocs.size()-1].timestamp.toInstant().toEpochMilli()-sortedDocs[0].timestamp.toInstant().toEpochMilli()); return uptimes;",
		"reduce_script": {"source": "def sum = 0.0; def range = 0.0; for (a in states) { for (int i=0;i<a.size()-1;i++){ sum+=a[i]; } range = a[a.size()-1]; } return sum/range;"}
		}
		}
		}
		}
		}
		}`, startMils, endMils)
	res, err := service.escli.Query(query)
	var response dto.Response
	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		log.Error(fmt.Sprintf("Error decoding the response: %s", err))
		return []dto.ServerUptime{}
	}
	return response.Aggregtions.Server.Buckets
}
