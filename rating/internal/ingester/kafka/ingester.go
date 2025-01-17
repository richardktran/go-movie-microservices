package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/richardktran/go-movie-microservices/rating/pkg/model"
)

// Ingester is a Kafka consumer that reads messages from a Kafka topic and sends them to the rating service.
type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Ingester{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
				return
			default:
			}

			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Println("Consume error: ", err.Error())
				continue
			}

			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Unmarshal error: ", err.Error())
				continue
			}

			// send event to channel ch
			ch <- event
		}
	}()

	return ch, nil
}
