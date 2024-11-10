package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
)

func (h *Handler) AssignedOrder() message.NoPublishHandlerFunc {
	type Event struct {
		OrderID   uuid.UUID `json:"order_id"`
		CourierID uuid.UUID `json:"courier_id"`
	}

	return func(msg *message.Message) error {
		var e Event

		if err := h.Unmarshaler(msg.Payload, &e); err != nil {
			return err
		}

		_, err := h.orderRepository.Transact(msg.Context(), e.OrderID, func(o *order.Order) error {
			o.CourierID = e.CourierID
			return o.SetCourierStatus(order.CourierAssigned)
		})
		if err != nil {
			return err
		}

		return err
	}
}
