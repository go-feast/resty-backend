package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
	"log"
)

func (h *Handler) PaymentWaiting() message.NoPublishHandlerFunc {
	type Event struct {
		OrderID              uuid.UUID `json:"order_id"`
		PaymentTransactionID uuid.UUID `json:"payment_id"`
	}

	return func(msg *message.Message) error {
		var event Event
		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return err
		}

		log.Printf("[payment.waiting] Received event: order_id:%s, payment_id:%s", event.OrderID, event.PaymentTransactionID)

		_, err := h.orderRepository.Transact(msg.Context(), event.OrderID, func(o *order.Order) error {
			o.PaymentID = event.PaymentTransactionID
			return o.SetPaymentStatus(order.PaymentWaiting)
		})
		if err != nil {
			return err
		}

		return nil
	}
}
