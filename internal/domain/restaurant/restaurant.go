package restaurant

import (
	"context"
	"fmt"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
)

type (
	RestaurantRepository interface {
		CreateRestaurant(ctx context.Context, restaurant *Restaurant) error
		GetRestaurant(ctx context.Context, id uuid.UUID) (*Restaurant, error)
		Transact(ctx context.Context, id uuid.UUID, action func(restaurant *Restaurant) error) (*Restaurant, error)
	}

	OrderRepository interface {
		CreateOrder(ctx context.Context, order *Order) error
		GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
		Transact(ctx context.Context, id uuid.UUID, action func(order *Order) error) (*Order, error)
	}

	Restaurant struct {
		ID                uuid.UUID `json:"id"`
		Name              string    `json:"name"`
		LocationLatitude  float64   `json:"latitude"`
		LocationLongitude float64   `json:"longitude"`
		Meals             []Meal    `json:"meals"`
	}

	Order struct {
		ID               uuid.UUID        `json:"id"`
		OrderStatus      OrderStatus      `json:"order_status"`
		RestaurantStatus RestaurantStatus `json:"restaurant_status"`
		RestaurantID     uuid.UUID        `json:"restaurant_id"`
		Meals            []OrderMeal      `json:"meals"`
	}

	Meal struct {
		ID           uuid.UUID `json:"id"`
		RestaurantID uuid.UUID `json:"restaurant_id"`
		Name         string    `json:"name"`
	}

	OrderMeal struct {
		ID      uuid.UUID
		OrderID uuid.UUID
	}

	OrderStatus      string
	RestaurantStatus string
)

func (s RestaurantStatus) String() string { return string(s) }
func (s OrderStatus) String() string      { return string(s) }

func mapMeals(oid uuid.UUID, meals uuid.UUIDs) []OrderMeal {
	m := make([]OrderMeal, len(meals))
	for i, meal := range meals {
		m[i] = OrderMeal{
			ID:      meal,
			OrderID: oid,
		}
	}

	return m
}

func NewOrder(
	ID uuid.UUID,
	restaurantID uuid.UUID,
	meals uuid.UUIDs,
) *Order {
	return &Order{
		ID:               ID,
		OrderStatus:      OrderCreated,
		RestaurantStatus: ReceivedOrder,
		RestaurantID:     restaurantID,
		Meals:            mapMeals(ID, meals),
	}
}

func NewRestaurant(
	name string,
	location geo.Location,
	meals []string,
) *Restaurant {
	m := make([]Meal, len(meals))

	id := uuid.New()
	for i, meal := range meals {
		m[i] = Meal{
			ID:           uuid.New(),
			RestaurantID: id,
			Name:         meal,
		}
	}

	return &Restaurant{
		ID:                id,
		Name:              name,
		LocationLongitude: location.Longitude,
		LocationLatitude:  location.Latitude,
		Meals:             m,
	}
}

func (o *Order) SetOrderStatus(status OrderStatus) error {
	if status == OrderCreated {
		if o.OrderStatus != OrderCreated && o.OrderStatus != "" {
			return fmt.Errorf("order status: %s", o.OrderStatus)
		} else {
			o.OrderStatus = status
		}
	}
	if status == OrderCanceled {
		if o.OrderStatus != OrderCompleted {
			o.OrderStatus = OrderCanceled
		} else {
			return fmt.Errorf("order status: %s", o.OrderStatus)
		}
	}
	if status == OrderCompleted {
		if o.OrderStatus != OrderCanceled {
			o.OrderStatus = OrderCompleted
		} else {
			return fmt.Errorf("order status: %s", o.OrderStatus)
		}
	}

	return nil
}

func (o *Order) IsCanceledOrCompleted() bool {
	return o.OrderStatus == OrderCanceled || o.OrderStatus == OrderCompleted
}

func (o *Order) SetRestaurantStatus(status RestaurantStatus) error {
	if status == ReceivedOrder {
		if o.RestaurantStatus != ReceivedOrder && o.RestaurantStatus != "" {
			return fmt.Errorf("restaurant status: %s", o.RestaurantStatus)
		} else {
			o.RestaurantStatus = status
		}
	}
	if status == PreparingOrder {
		if o.RestaurantStatus != PreparedOrder && !o.IsCanceledOrCompleted() {
			o.RestaurantStatus = PreparingOrder
		} else {
			return fmt.Errorf("restaurant status: %s", o.RestaurantStatus)
		}
	}
	if status == PreparedOrder && !o.IsCanceledOrCompleted() {
		o.RestaurantStatus = PreparedOrder
	}

	return nil
}

const (
	ReceivedOrder  RestaurantStatus = "restaurant.order.received"
	PreparingOrder RestaurantStatus = "restaurant.order.preparing"
	PreparedOrder  RestaurantStatus = "restaurant.order.prepared"
)

const (
	OrderCreated   OrderStatus = "order.created"
	OrderCanceled  OrderStatus = "order.canceled"
	OrderCompleted OrderStatus = "order.completed"
)
