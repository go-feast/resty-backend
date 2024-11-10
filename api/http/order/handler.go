package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/order"
)

type Handler struct {
	orderRepository order.OrderRepository
	publisher       message.Publisher
	Marshaler       func(interface{}) ([]byte, error)
}

func NewHandler(
	orderRepository order.OrderRepository,
	publisher message.Publisher,
	marshaler func(interface{}) ([]byte, error),
) *Handler {
	return &Handler{
		orderRepository: orderRepository,
		publisher:       publisher,
		Marshaler:       marshaler,
	}
}
