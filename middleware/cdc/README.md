# CDC

```yml
# docker-compose.yml
services:
  zookeeper:
    image: wurstmeister/zookeeper   ## 镜像
    container_name: zookeeper
    ports:
      - "2181:2181"                 ## 对外暴露的端口号
    volumes:
      - ./conf/zoo.cfg:/conf/zoo.cfg
      - ./data/zookeeper:/data
      - ./logs/zookeeper:/datalog
    restart: always
  kafka:
    image: wurstmeister/kafka       ## 镜像
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: INSIDE://10.66.0.11:9092 #注意,这里不能设置成localhost和127.0.0.1
      KAFKA_LISTENERS: INSIDE://0.0.0.0:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT
      KAFKA_LISTENER_NAME_SELECTOR: INSIDE
      KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE
    volumes:
      - ./logs/kafka:/opt/kafka/logs
    restart: always
  kafka-ui:
    image: provectuslabs/kafka-ui
    container_name: kafka-ui
    ports:
       - "9080:8080"
    depends_on:
      - kafka
    environment:
      DYNAMIC_CONFIG_ENABLED: 'true'
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAP_SERVERS: kafka:9092  # 使用 Kafka 服务的内部网络地址
    restart: always
  kafdrop: # kafka的图形化界面工具
    image: obsidiandynamics/kafdrop
    container_name: kafdrop
    ports:
      - "9000:9000"
    depends_on:
      - kafka
    environment:
      - KAFKA_BROKERCONNECT=kafka:9092
      - SERVER_SERVLET_CONTEXTPATH=/
    restart: always
  # debezium
  debezium-connector:
    image: debezium/connect:2.7.3.Final
    container_name: debezium-connector
    #depends_on:
    #  - zookeeper
    #  - kafka
    ports:
      - "8083:8083"
    environment:
      BOOTSTRAP_SERVERS: 10.66.0.11:9092
      GROUP_ID: debezium-group
      CONFIG_STORAGE_TOPIC: connect-configs
      OFFSET_STORAGE_TOPIC: connect-offsets
      STATUS_STORAGE_TOPIC: connect-status
  kafka-connect-ui:
    image: landoop/kafka-connect-ui:0.9.7
    restart: always
    container_name: kafka-connect-ui
    links:
      - debezium-connector
    ports:
      - "8008:8000"
    environment:
      - CONNECT_URL=http://debezium-connector:8083
```

```shell
curl -X POST http://10.66.0.11:8008/api/kafka-connect-1/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "go-admin-sys_login_log",  
    "config": {
      "connector.class": "io.debezium.connector.mysql.MySqlConnector",
	  "database.user": "root",
	  "topic.creation.default.partitions": "5",
	  "database.server.id": "1",
	  "tasks.max": "1",
	  "database.server.name": "mysql-server",
	  "schema.history.internal.kafka.bootstrap.servers": "10.66.0.11:9092",
	  "database.port": "3306",
	  "include.schema.changes": "true",
	  "topic.prefix": "v_debezium",
	  "schema.history.internal.kafka.topic": "go-admin-history",
	  "database.hostname": "10.66.0.11",
	  "database.password": "123456",
	  "topic.creation.default.replication.factor": "1",
	  "table.include.list": "go-admin.sys_login_log",
	  "database.include.list": "go-admin"
    }
  }'
```

 