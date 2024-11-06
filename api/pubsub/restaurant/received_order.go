package restaurant

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	internalmessage "github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
)

func (h *Handler) ReceivedOrder() message.HandlerFunc {
	type Event struct {
		OrderID      uuid.UUID  `json:"order_id"`
		RestaurantID uuid.UUID  `json:"restaurant_id"`
		Meals        uuid.UUIDs `json:"meals"`
	}

	return func(msg *message.Message) ([]*message.Message, error) {
		var event Event

		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return nil, err
		}

		o := restaurant.NewOrder(event.OrderID, event.RestaurantID, event.Meals)

		if err := h.orderRepository.CreateOrder(msg.Context(), o); err != nil {
			return nil, err
		}

		msg = internalmessage.NewMessage(internalmessage.Event{
			"order_id": event.OrderID,
		}, h.Marshaler)

		return []*message.Message{msg}, nil
	}
}
