package cdc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
)

type KafkaMessage struct {
	Before json.RawMessage `json:"before"`
	After  json.RawMessage `json:"after"`
	Source struct {
		Version   string `json:"version"`
		Connector string `json:"connector"`
		Name      string `json:"name"`
		ServerID  int64  `json:"server_id"`
		TSMS      int64  `json:"ts_ms"`
		Table     string `json:"table"`
		Database  string `json:"database"`
	} `json:"source"`
}

func sarama_start() {
	// Kafka消费者配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	/* partitionConsumer, err := consumer.ConsumePartition(
		"mysql-server.your_database.your_table",
		sarama.OffsetNewest,
		nil,
	) */
	partitionConsumer, err := consumer.ConsumePartition(
		"mysql-server.your_database.your_table",
		int32(sarama.OffsetNewest),
		0,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionConsumer.Close()

	// Elasticsearch客户端
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	// 消费消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var kmsg KafkaMessage
			if err := json.Unmarshal(msg.Value, &kmsg); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			if len(kmsg.After) == 0 {
				continue // 只处理INSERT/UPDATE事件
			}

			// 解析数据
			var data map[string]interface{}
			if err := json.Unmarshal(kmsg.After, &data); err != nil {
				log.Printf("Error unmarshalling data: %v", err)
				continue
			}

			// 添加元数据
			data["_id"] = kmsg.Source.ServerID
			data["_source"] = data["_source"].(map[string]interface{})
			delete(data, "_source")

			fmt.Printf("%v", es.Index)
			// 写入ES TODO
			/* res, err := es.Index(
				es.Index.WithIndexName("mysql_data"),
				es.Reader(strings.NewReader(fmt.Sprintf("%v", data))),
			)
			if err != nil {
				log.Printf("ES write error: %v", err)
			} else {
				log.Printf("Indexed document: %s", res.String())
			} */
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

func consumer_start() {
	// Kafka消费者配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	/* partitionConsumer, err := consumer.ConsumePartition(
		"mysql-server.your_database.your_table",
		sarama.OffsetNewest,
		nil,
	) */
	partitionConsumer, err := consumer.ConsumePartition(
		"v_debezium.*",
		int32(sarama.OffsetNewest),
		0,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionConsumer.Close()

	// 消费消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var kmsg KafkaMessage
			if err := json.Unmarshal(msg.Value, &kmsg); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			if len(kmsg.After) == 0 {
				continue // 只处理INSERT/UPDATE事件
			}

			// 解析数据
			var data map[string]interface{}
			if err := json.Unmarshal(kmsg.After, &data); err != nil {
				log.Printf("Error unmarshalling data: %v", err)
				continue
			}

			// 添加元数据
			data["_id"] = kmsg.Source.ServerID
			data["_source"] = data["_source"].(map[string]interface{})
			delete(data, "_source")
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}
