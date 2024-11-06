package restaurant

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	"github.com/google/uuid"
)

func (h *Handler) OrderCompleted() message.NoPublishHandlerFunc {
	type Event struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	return func(msg *message.Message) error {
		var event Event

		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return err
		}

		_, err := h.orderRepository.Transact(msg.Context(), event.OrderID, func(order *restaurant.Order) error {
			return order.SetOrderStatus(restaurant.OrderCompleted)
		})
		if err != nil {
			return err
		}

		return nil
	}
}
