package order

import (
	"github.com/go-chi/render"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/shared/saver"
	"github.com/go-feast/resty-backend/internal/http/httpstatus"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

func (h *Handler) TakeOrder() func(http.ResponseWriter, *http.Request) {
	return takeOrder(h.tracer, h.saverService, h.restaurantRepository)
}

func takeOrder(tracer trace.Tracer, s saver.TransactionalOutbox[*order.Order], restaurantRepository order.RestaurantRepository) func(w http.ResponseWriter, r *http.Request) {
	type TakeOrderRequest struct {
		CustomerID   string   `json:"customer_id"`
		RestaurantID string   `json:"restaurant_id"`
		Meals        []string `json:"meals"`
	}

	type TakeOrderResponse struct { //nolint:govet
		OrderID   string    `json:"order_id"`
		Timestamp time.Time `json:"timestamp"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx, span        = tracer.Start(r.Context(), "take order")
			takeOrderRequest = &TakeOrderRequest{}
		)

		defer span.End()

		if err := render.DecodeJSON(r.Body, takeOrderRequest); err != nil {
			httpstatus.BadRequest(ctx, w, err)
			return
		}

		rid := uuid.MustParse(takeOrderRequest.RestaurantID)
		restaurant, err := restaurantRepository.GetByID(ctx, rid)
		if err != nil {
			httpstatus.BadRequest(ctx, w, err)
			return
		}

		o, err := order.NewOrder(
			takeOrderRequest.RestaurantID,
			takeOrderRequest.CustomerID,
			takeOrderRequest.Meals,
			restaurant.Location.Latitude(),
			restaurant.Location.Longitude(),
		)
		if err != nil {
			httpstatus.BadRequest(ctx, w, err)
			return
		}

		err = s.Save(ctx, o)
		if err != nil {
			// should be bad request or internal service error
			httpstatus.InternalServerError(ctx, w, err)
			return
		}

		span.AddEvent("order created")

		response := TakeOrderResponse{
			OrderID:   o.ID().String(),
			Timestamp: o.CreateAt(),
		}

		httpstatus.Created(w, response)
	}
}
