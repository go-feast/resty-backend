package order

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/internal/http/httpstatus"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (h *Handler) GetOrder() func(w http.ResponseWriter, r *http.Request) {
	return getOrder(h.repository)
}

func getOrder(repository order.Repository) func(w http.ResponseWriter, r *http.Request) {
	type GetOrderResponse struct { //nolint:govet
		ID            uuid.UUID `json:"id"`
		State         string    `json:"state"`
		TransactionID string    `json:"transaction_id"`
		CourierID     string    `json:"courier_id"`
		Timestamp     time.Time `json:"timestamp"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "uuid"))
		if err != nil {
			httpstatus.BadRequest(ctx, w, errors.Wrap(err, "invalid order id"))
			return
		}

		o, err := repository.Get(ctx, id)
		if err != nil {
			httpstatus.InternalServerError(ctx, w, errors.Wrap(err, "failed to get order"))
			return
		}

		response := GetOrderResponse{
			ID:            o.ID(),
			State:         o.State().String(),
			TransactionID: o.TransactionID().String(),
			CourierID:     o.CourierID().String(),
			Timestamp:     time.Now(),
		}

		render.JSON(w, r, response)
	}
}
