package event

import (
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
)

// JSONEventOrderCreated provides JSON representation of Order.
type JSONEventOrderCreated struct {
	OrderID      string           `json:"order_id"`
	CustomerID   string           `json:"customer_id"`
	RestaurantID string           `json:"restaurant_id"`
	Meals        []string         `json:"meals"`
	Destination  geo.JSONLocation `json:"destination"`
}

type JSONOrderFinished struct {
	OrderID uuid.UUID `json:"order_id"`
}

type JSONOrderCooking struct {
	OrderID uuid.UUID `json:"order_id"`
}

type JSONEventOrderPaid struct {
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
}

type JSONWaitingForCourier struct {
	OrderID uuid.UUID `json:"order_id"`
}

type JSONCourierTook struct {
	OrderID   uuid.UUID `json:"order_id"`
	CourierID uuid.UUID `json:"courier_id"`
}

type JSONDelivering struct {
	OrderID uuid.UUID `json:"order_id"`
}

type JSONDelivered struct {
	OrderID uuid.UUID `json:"order_id"`
}

type JSONCanceled struct { //nolint:govet
	OrderID uuid.UUID `json:"order_id"`
	Reason  string    `json:"reason"`
}

type JSONClosed struct {
	OrderID uuid.UUID `json:"order_id"`
}
