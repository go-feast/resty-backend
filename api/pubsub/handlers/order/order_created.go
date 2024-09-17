package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/pubsub"
	"github.com/rs/zerolog/log"
)

func (h *Handler) OrderCreated(msg *message.Message) ([]*message.Message, error) {
	_, span := pubsub.SpanFromMessage(
		msg,
		"consumer.order.created",
		"order.created handler",
		nil,
	)
	defer span.End()

	log.Info().Str("msg-id", msg.UUID).Msg("Received message from topic OrderCreated")

	return []*message.Message{msg}, nil
}
