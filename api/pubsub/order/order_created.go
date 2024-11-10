package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
)

func (h *Handler) OrderCreated() message.NoPublishHandlerFunc {
	type Order struct {
		ID string `json:"order_id"`
	}
	return func(msg *message.Message) error {
		var order Order

		if err := h.Unmarshaler(msg.Payload, &order); err != nil {
			return err
		}

		log.Printf("Order [%s] created", order.ID)

		return nil
	}
}
