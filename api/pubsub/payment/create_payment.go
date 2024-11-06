package payment

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/payment"
	internalmsg "github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
)

// should listen on order.created and post to payment.waiting
func (h *Handler) CreatePayment() message.HandlerFunc {
	type Event struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	return func(msg *message.Message) ([]*message.Message, error) {
		var event Event
		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return nil, err
		}

		p := payment.NewPayment(event.OrderID)

		if err := h.paymentRepository.Create(msg.Context(), p); err != nil {
			return nil, err
		}

		msg = internalmsg.NewMessage(internalmsg.Event{"order_id": event.OrderID}, h.Marshaler)
		return []*message.Message{msg}, nil
	}
}
