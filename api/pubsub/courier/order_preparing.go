package courier

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/courier"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"time"
)

func (h *Handler) OrderPreparing() message.NoPublishHandlerFunc {
	type Event struct {
		OrderID            uuid.UUID    `json:"order_id"`
		RestaurantLocation geo.Location `json:"restaurant_location"`
		Destination        geo.Location `json:"destination"`
		CreatedAt          time.Time    `json:"created_at"`
	}

	return func(msg *message.Message) error {
		var e Event
		if err := h.Unmarshaler(msg.Payload, &e); err != nil {
			return err
		}

		order := courier.NewOrder(
			e.OrderID,
			e.RestaurantLocation,
			e.Destination,
			e.CreatedAt,
		)

		ctx, cancel := context.WithTimeout(msg.Context(), 5*time.Second)
		defer cancel()

		err := h.courierRepository.CreateOrder(ctx, order)
		if err != nil {
			return err
		}

		return nil
	}
}
