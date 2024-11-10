package order

import "github.com/go-feast/resty-backend/internal/domain/order"

type Handler struct {
	orderRepository order.OrderRepository
	Unmarshaler     func([]byte, interface{}) error
	Marshaler       func(interface{}) ([]byte, error)
}

func NewHandler(
	orderRepository order.OrderRepository,
	unmarshaler func([]byte, interface{}) error,
	marshaler func(interface{}) ([]byte, error),
) *Handler {
	return &Handler{
		orderRepository: orderRepository,
		Unmarshaler:     unmarshaler,
		Marshaler:       marshaler,
	}
}
