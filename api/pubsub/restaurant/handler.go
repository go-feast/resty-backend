package restaurant

import (
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
)

type Handler struct {
	orderRepository restaurant.OrderRepository
	Unmarshaler     func([]byte, interface{}) error
	Marshaler       func(interface{}) ([]byte, error)
}

func NewHandler(
	orderRepository restaurant.OrderRepository,
	unmarshaler func([]byte, interface{}) error,
	marshaler func(interface{}) ([]byte, error),
) *Handler {
	return &Handler{
		orderRepository: orderRepository,
		Unmarshaler:     unmarshaler,
		Marshaler:       marshaler,
	}
}
