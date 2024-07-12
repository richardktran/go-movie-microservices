package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/richardktran/go-movie-microservices/rating/pkg/model"
)

func main() {
	fmt.Println("Creating a Kafka producer...")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})

	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratingsdata.json"
	fmt.Printf("Reading data from %s...\n", fileName)
	ratingEvents, err := readRatingEvents(fileName)

	if err != nil {
		panic(err)
	}

	topic := "ratings"

	if err := produceRatingEvents(producer, topic, ratingEvents); err != nil {
		panic(err)
	}

	var timeout = 10 * time.Second
	fmt.Printf("Waiting %s until all events get produced...\n", timeout.String())
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}

func produceRatingEvents(producer *kafka.Producer, topic string, ratingEvents []model.RatingEvent) error {
	for _, event := range ratingEvents {
		encodedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}
		fmt.Printf("Sending event: %s\n", encodedEvent)

		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(encodedEvent),
		}, nil); err != nil {
			return err
		}
	}

	return nil
}
