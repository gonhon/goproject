package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/segmentio/kafka-go"
)

// ChangeEvent 表示 Debezium 变更事件结构
type ChangeEvent struct {
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
	Source struct {
		Database string `json:"db"`
		Table    string `json:"table"`
	} `json:"source"`
	Op string `json:"op"` // "c"=create, "u"=update, "d"=delete
}

// ESClient 封装 Elasticsearch 操作
type ESClient struct {
	client *elasticsearch.Client
}

func NewESClient(addresses []string) (*ESClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ESClient{client: client}, nil
}

// IndexDocument 索引文档到指定索引
func (es *ESClient) IndexDocument(index string, id string, doc map[string]interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

// UpdateDocument 更新 Elasticsearch 中的文档
func (es *ESClient) UpdateDocument(index string, id string, doc map[string]interface{}) error {
	body, err := json.Marshal(map[string]interface{}{"doc": doc})
	if err != nil {
		return fmt.Errorf("error marshaling update document: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

// DeleteDocument 从 Elasticsearch 删除文档
func (es *ESClient) DeleteDocument(index string, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

// getDocumentID 从变更事件中提取文档ID
func getDocumentID(doc map[string]interface{}) (string, error) {
	// 根据你的主键字段调整
	if id, ok := doc["id"].(string); ok {
		return id, nil
	}
	if id, ok := doc["ID"].(string); ok {
		return id, nil
	}
	if id, ok := doc["Id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("could not determine document ID from fields: %v", doc)
}

// processEvent 处理单个变更事件
func processEvent(es *ESClient, event ChangeEvent) error {
	// 构建 Elasticsearch 索引名称 (db_table 格式)
	indexName := fmt.Sprintf("%s_%s", strings.ToLower(event.Source.Database), strings.ToLower(event.Source.Table))

	switch event.Op {
	case "c": // 创建文档
		id, err := getDocumentID(event.After)
		if err != nil {
			return fmt.Errorf("error getting document ID for create: %w", err)
		}
		return es.IndexDocument(indexName, id, event.After)
	case "u": // 更新文档
		id, err := getDocumentID(event.After)
		if err != nil {
			return fmt.Errorf("error getting document ID for update: %w", err)
		}
		return es.UpdateDocument(indexName, id, event.After)
	case "d": // 删除文档
		id, err := getDocumentID(event.Before)
		if err != nil {
			return fmt.Errorf("error getting document ID for delete: %w", err)
		}
		return es.DeleteDocument(indexName, id)
	default:
		return fmt.Errorf("unknown operation type: %s", event.Op)
	}
}

func segmentio_start() {
	// 初始化 Elasticsearch 客户端
	esClient, err := NewESClient([]string{"http://localhost:9200"})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

	// 初始化 Kafka 消费者
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "es-sync-group",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		// 监听所有表变更的 topic 模式
		Topic: "v_debezium.*",
		// GroupTopics: []string{"dbserver1.your_database.table1", "dbserver1.your_database.table2"},
	})

	log.Println("Starting consumer...")
	for {
		// 读取消息
		msg, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// 解析变更事件
		var event ChangeEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling event from topic %s: %v", msg.Topic, err)
			continue
		}

		// 处理事件
		if err := processEvent(esClient, event); err != nil {
			log.Printf("Error processing event from topic %s: %v", msg.Topic, err)
			// 这里可以添加重试逻辑或死信队列处理
			continue
		}

		log.Printf("Successfully processed %s operation for table %s.%s",
			map[string]string{"c": "create", "u": "update", "d": "delete"}[event.Op],
			event.Source.Database,
			event.Source.Table)
	}
}

// 在 main 函数初始化后创建索引
// table1Mapping := `{
// 	"mappings": {
// 		"properties": {
// 			"id": {"type": "keyword"},
// 			"name": {"type": "text"},
// 			"created_at": {"type": "date"}
// 		}
// 	}
// }`

// if err := esClient.CreateIndexIfNotExists("your_database_table1", table1Mapping); err != nil {
// 	log.Fatalf("Error creating index for table1: %v", err)
// }

// CreateIndexIfNotExists 创建索引（如果不存在）
func (es *ESClient) CreateIndexIfNotExists(index string, mapping string) error {
	existsReq := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := existsReq.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil // 索引已存在
	}

	createReq := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(mapping),
	}

	res, err = createReq.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

// transformFields 对字段进行转换处理
func transformFields(source map[string]interface{}, table string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range source {
		// 统一处理字段名
		newKey := strings.ToLower(k)

		// 表特定字段处理
		switch table {
		case "users":
			if newKey == "birthdate" {
				if dateStr, ok := v.(string); ok {
					if t, err := time.Parse("2006-01-02", dateStr); err == nil {
						result[newKey] = t.Format(time.RFC3339)
						continue
					}
				}
			}
		case "products":
			if newKey == "price" {
				if priceStr, ok := v.(string); ok {
					result[newKey+"_cents"] = strings.Replace(priceStr, ".", "", -1)
				}
			}
		}

		// 默认处理
		result[newKey] = v
	}

	return result
}

// BulkProcessor 批量处理变更事件
type BulkProcessor struct {
	esClient    *ESClient
	batchSize   int
	maxInterval time.Duration
	queue       chan ChangeEvent
}

func NewBulkProcessor(esClient *ESClient, batchSize int, maxInterval time.Duration) *BulkProcessor {
	return &BulkProcessor{
		esClient:    esClient,
		batchSize:   batchSize,
		maxInterval: maxInterval,
		queue:       make(chan ChangeEvent, batchSize*2),
	}
}

func (bp *BulkProcessor) Start() {
	go func() {
		var batch []ChangeEvent
		ticker := time.NewTicker(bp.maxInterval)
		defer ticker.Stop()

		for {
			select {
			case event := <-bp.queue:
				batch = append(batch, event)
				if len(batch) >= bp.batchSize {
					bp.processBatch(batch)
					batch = nil
					ticker.Reset(bp.maxInterval)
				}
			case <-ticker.C:
				if len(batch) > 0 {
					bp.processBatch(batch)
					batch = nil
				}
			}
		}
	}()
}

func (bp *BulkProcessor) AddEvent(event ChangeEvent) {
	bp.queue <- event
}

func (bp *BulkProcessor) processBatch(batch []ChangeEvent) {
	var buf strings.Builder

	for _, event := range batch {
		indexName := fmt.Sprintf("%s_%s", strings.ToLower(event.Source.Database), strings.ToLower(event.Source.Table))
		id, err := getDocumentID(event.After)
		if err != nil {
			log.Printf("Error getting document ID in batch: %v", err)
			continue
		}

		switch event.Op {
		case "c", "u":
			// { "index" : { "_index" : "test", "_id" : "1" } }
			meta := map[string]interface{}{
				"index": map[string]interface{}{
					"_index": indexName,
					"_id":    id,
				},
			}

			metaJSON, _ := json.Marshal(meta)
			buf.Write(metaJSON)
			buf.WriteByte('\n')

			docJSON, _ := json.Marshal(event.After)
			buf.Write(docJSON)
			buf.WriteByte('\n')
		case "d":
			// { "delete" : { "_index" : "test", "_id" : "2" } }
			meta := map[string]interface{}{
				"delete": map[string]interface{}{
					"_index": indexName,
					"_id":    id,
				},
			}

			metaJSON, _ := json.Marshal(meta)
			buf.Write(metaJSON)
			buf.WriteByte('\n')
		}
	}

	if buf.Len() > 0 {
		res, err := bp.esClient.client.Bulk(
			strings.NewReader(buf.String()),
			bp.esClient.client.Bulk.WithRefresh("true"),
		)
		if err != nil {
			log.Printf("Error executing bulk request: %v", err)
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("Error response from bulk request: %s", res.String())
		}
	}
}
