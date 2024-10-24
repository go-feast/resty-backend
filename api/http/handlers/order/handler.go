package order

import (
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/shared/saver"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	tracer trace.Tracer

	saverService saver.TransactionalOutbox[*order.Order]

	// metrics

	// repositories eg.
	orderRepository      order.OrderRepository
	restaurantRepository order.RestaurantRepository
}

func NewHandler(
	tracer trace.Tracer,
	repository order.OrderRepository,
	saverService saver.TransactionalOutbox[*order.Order],
	restaurantRepository order.RestaurantRepository,
) *Handler {
	return &Handler{
		tracer:               tracer,
		orderRepository:      repository,
		saverService:         saverService,
		restaurantRepository: restaurantRepository,
	}
}
