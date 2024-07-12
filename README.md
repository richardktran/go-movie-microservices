## Run Hashicorp Consul in Docker
```bash
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul:latest agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```

## Run Kafka in Docker
Build the image and run the container
```bash
cd cmd/ratingingester
docker-compose up -d
```
Access Kafka
```bash
docker exec -it kafka /bin/sh
cd /opt/kafka_2.13-2.8.1/bin
```
Create Topic
```bash
kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic ratings
```
Send a message to the topic (Open a new terminal)
```bash
kafka-console-producer.sh --bootstrap-server localhost:9092 --topic ratings
<Enter your message>
```
Consume messages from the topic (Open a new terminal)
```bash
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic ratings --from-beginning --max-messages 20
```

### To run example
Run producer
```bash
cd cmd/ratingingester
go run main.go
```
Run consumer
```bash
cd rating/cmd/consumer
go run main.go
```
