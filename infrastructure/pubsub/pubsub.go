package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/config"
	"log"
	"time"
)

func NewKafkaSub(consumerGroup string) message.Subscriber {
	cfg := kafka.SubscriberConfig{
		Brokers:             config.Kafka(),
		Unmarshaler:         kafka.DefaultMarshaler{},
		ConsumerGroup:       consumerGroup,
		ReconnectRetrySleep: 1 * time.Second,
	}

	subscriber, err := kafka.NewSubscriber(cfg, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create kafka subscriber: %v", err)
	}

	return subscriber
}

func NewKafkaPub() message.Publisher {
	cfg := kafka.PublisherConfig{
		Brokers:   config.Kafka(),
		Marshaler: kafka.DefaultMarshaler{},
	}

	publisher, err := kafka.NewPublisher(cfg, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatalf("Failed to create kafka publisher: %v", err)
	}

	return publisher
}
