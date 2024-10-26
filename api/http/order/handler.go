package order

import "github.com/go-feast/resty-backend/internal/domain/order"

type Handler struct {
	orderRepository order.OrderRepository
}

func NewHandler(orderRepository order.OrderRepository) *Handler {
	return &Handler{orderRepository: orderRepository}
}
