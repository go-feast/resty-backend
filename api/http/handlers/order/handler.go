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
	repository order.Repository
}

func NewHandler(
	tracer trace.Tracer,
	repository order.Repository,
	saverService saver.TransactionalOutbox[*order.Order],
) *Handler {
	return &Handler{
		tracer:       tracer,
		repository:   repository,
		saverService: saverService,
	}
}
