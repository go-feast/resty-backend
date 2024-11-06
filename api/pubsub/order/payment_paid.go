package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
)

func (h *Handler) PaymentPaid() message.NoPublishHandlerFunc {
	type Event struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	return func(msg *message.Message) error {
		var event Event
		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return err
		}

		_, err := h.orderRepository.Transact(msg.Context(), event.OrderID, func(o *order.Order) error {
			o.PaymentStatus = order.PaymentPaid
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}
}
