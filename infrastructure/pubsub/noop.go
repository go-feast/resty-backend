package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
)

type NopPublisher struct{}

func (n *NopPublisher) Publish(topic string, messages ...*message.Message) error {
	log.Printf("[NOOP PUBLISHER] [%s] Message published: %s", topic, messages[0].UUID)
	return nil
}

func (n *NopPublisher) Close() error { return nil }
