package order

import (
	"github.com/go-feast/resty-backend/domain/order"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	unmarshaler func([]byte, any) error
	tracer      trace.Tracer
	repository  order.OrderRepository
}

func NewHandler(
	unmarshaler func([]byte, any) error,
	tracer trace.Tracer,
	repository order.OrderRepository,
) *Handler {
	return &Handler{
		unmarshaler: unmarshaler,
		tracer:      tracer,
		repository:  repository,
	}
}
