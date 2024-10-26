package order

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
}

type Order struct {
	ID               uuid.UUID
	CourierID        uuid.UUID
	RestaurantID     uuid.UUID
	CustomerID       uuid.UUID
	DestinationLat   float64
	DestinationLng   float64
	Meals            []Meal `gorm:"foreignKey:OrderID"`
	TransactionID    uuid.UUID
	RestaurantStatus RestaurantStatus
	PaymentStatus    TransactionStatus
	CourierStatus    CourierStatus
	OrderStatus      OrderStatus
	CreatedAt        time.Time
	CancelledAt      time.Time
	DeletedAt        time.Time
}

type Destination struct {
	gorm.Model
	OrderID  uuid.UUID
	Location geo.Location
}

type Meal struct {
	ID      uuid.UUID
	OrderID uuid.UUID
}

func mapMeals(oid uuid.UUID, meals uuid.UUIDs) []Meal {
	m := make([]Meal, len(meals))
	for i, meal := range meals {
		m[i] = Meal{
			ID:      meal,
			OrderID: oid,
		}
	}

	return m
}

func NewOrder(customerID uuid.UUID, restaurantID uuid.UUID, meals uuid.UUIDs, destination geo.Location) *Order {
	id := uuid.New()
	return &Order{
		ID:               id,
		CourierID:        uuid.Nil,
		RestaurantID:     restaurantID,
		CustomerID:       customerID,
		DestinationLat:   destination.Latitude,
		DestinationLng:   destination.Longitude,
		Meals:            mapMeals(id, meals),
		TransactionID:    uuid.Nil,
		RestaurantStatus: "",
		PaymentStatus:    "",
		CourierStatus:    "",
		OrderStatus:      Created,
		CreatedAt:        time.Now(),
		CancelledAt:      time.Time{},
		DeletedAt:        time.Time{},
	}
}

type (
	TransactionStatus string
	RestaurantStatus  string
	OrderStatus       string
	CourierStatus     string
)

const (
	Created  OrderStatus = "order.created"
	Canceled OrderStatus = "order.cancelled"
	Closed   OrderStatus = "order.closed"
)

const (
	CourierTookOrder  = "courier.order.taken"
	CourierDelivering = "courier.delivering"
	CourierDelivered  = "courier.delivered"
)

const (
	RestaurantReceivedOrder  RestaurantStatus = "restaurant.order.received"
	RestaurantPreparingOrder RestaurantStatus = "restaurant.order.preparing"
	RestaurantPreparedOrder  RestaurantStatus = "restaurant.order.prepared"
)

const (
	PaymentWaiting  TransactionStatus = "payment.waiting"
	PaymentPaid     TransactionStatus = "payment.paid"
	PaymentCanceled TransactionStatus = "payment.canceled"
)
