package order

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID               uuid.UUID
	CourierID        uuid.UUID
	RestaurantID     uuid.UUID
	Meals            uuid.UUIDs
	TransactionID    uuid.UUID
	RestaurantStatus RestaurantStatus
	PaymentStatus    TransactionStatus
	CourierStatus    CourierStatus
	OrderStatus      OrderStatus
	CreatedAt        time.Time
	CancelledAt      time.Time
	DeletedAt        time.Time
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
