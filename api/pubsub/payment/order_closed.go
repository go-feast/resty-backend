package payment

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/payment"
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

		if _, err := h.paymentRepository.Transact(msg.Context(), event.OrderID, func(p *payment.Payment) error {
			return p.SetOrderStatus(payment.OrderCompleted)
		}); err != nil {
			return err
		}

		return nil
	}
}
