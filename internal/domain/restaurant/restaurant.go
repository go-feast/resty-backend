package restaurant

import (
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
)

type (
	Restaurant struct {
		ID       uuid.UUID
		Name     string
		Location geo.Location
		Meals    []Meal
	}

	Order struct {
		ID           uuid.UUID
		Status       OrderStatus
		RestaurantID uuid.UUID
		Meals        uuid.UUID
	}

	Meal struct {
		ID   uuid.UUID
		Name string
	}

	OrderStatus string
)

const (
	ReceivedOrder  OrderStatus = "restaurant.order.received"
	PreparingOrder OrderStatus = "restaurant.order.preparing"
	PreparedOrder  OrderStatus = "restaurant.order.prepared"
)
