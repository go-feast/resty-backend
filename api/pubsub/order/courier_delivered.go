package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *Handler) CourierDelivered() message.HandlerFunc {
	type Event struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	return func(msg *message.Message) ([]*message.Message, error) {
		var event Event
		if err := h.Unmarshaler(msg.Payload, &event); err != nil {
			return nil, err
		}

		_, err := h.orderRepository.Transact(msg.Context(), event.OrderID, func(o *order.Order) error {

			if err := o.SetCourierStatus(order.CourierDelivered); err != nil {
				return errors.Wrap(err, "failed to set courier status")
			}

			if err := o.SetOrderStatus(order.Completed); err != nil {
				return errors.Wrap(err, "failed to set order status")
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return []*message.Message{msg}, nil
	}
}
