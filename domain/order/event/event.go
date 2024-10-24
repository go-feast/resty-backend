package event

import (
	"github.com/go-feast/resty-backend/domain/shared/geo"
)

// Type provides methods for converting Order for different marshaling strategies.
type Type struct {
	OrderID       string
	CustomerID    string
	RestaurantID  string
	TransactionID string
	Meals         []string
	Destination   geo.Location
}

// JSONEventOrderCreated converts Type to JSONEventOrderCreated.
func (t *Type) JSONEventOrderCreated() JSONEventOrderCreated {
	return JSONEventOrderCreated{
		OrderID:      t.OrderID,
		CustomerID:   t.CustomerID,
		RestaurantID: t.RestaurantID,
		Meals:        t.Meals,
		Destination:  t.Destination.ToJSON(),
	}
}
