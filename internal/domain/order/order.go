package order

import (
	"context"
	"fmt"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	Transact(ctx context.Context, id uuid.UUID, action func(o *Order) error) (*Order, error)
}

type Order struct {
	ID               uuid.UUID
	CourierID        uuid.UUID
	RestaurantID     uuid.UUID
	CustomerID       uuid.UUID
	PaymentID        uuid.UUID
	DestinationLat   float64
	DestinationLng   float64
	Meals            []Meal `gorm:"foreignKey:OrderID"`
	RestaurantStatus RestaurantStatus
	PaymentStatus    TransactionStatus
	CourierStatus    CourierStatus
	OrderStatus      OrderStatus
	CreatedAt        time.Time
	CancelledAt      time.Time
	DeletedAt        time.Time
}

func (o *Order) SetOrderStatus(status OrderStatus) error {
	if o.OrderStatus == "" {
		o.OrderStatus = status
	}

	if status == Created {
		if o.OrderStatus != Created {
			return fmt.Errorf("order satus: %s", o.OrderStatus)
		}
	}
	if status == Canceled {
		if o.OrderStatus != Completed {
			o.OrderStatus = Canceled
		} else {
			return fmt.Errorf("order satus: %s", o.OrderStatus)
		}
	}
	if status == Completed {
		o.OrderStatus = Completed
	}

	return nil
}

func (o *Order) SetPaymentStatus(status TransactionStatus) error {
	if status == PaymentWaiting {
		if o.PaymentStatus != PaymentWaiting && o.PaymentStatus != "" {
			return fmt.Errorf("payment satus: %s", o.PaymentStatus)
		} else {
			o.PaymentStatus = PaymentWaiting
		}
	}
	if status == PaymentCanceled {
		if o.PaymentStatus != PaymentPaid {
			o.PaymentStatus = PaymentCanceled
		} else {
			return fmt.Errorf("payment satus: %s", o.PaymentStatus)
		}
	}
	if status == PaymentPaid {
		if o.PaymentStatus != PaymentCanceled {
			o.PaymentStatus = PaymentPaid
		} else {
			return fmt.Errorf("payment satus: %s", o.PaymentStatus)
		}
	}

	return nil
}

func (o *Order) SetRestaurantStatus(status RestaurantStatus) error {
	if status == RestaurantReceivedOrder {
		if o.RestaurantStatus != RestaurantReceivedOrder && o.RestaurantStatus != "" {
			return fmt.Errorf("restaurant satus: %s", o.RestaurantStatus)
		} else {
			o.RestaurantStatus = RestaurantReceivedOrder
		}
	}
	if status == RestaurantPreparingOrder {
		if o.RestaurantStatus != RestaurantPreparedOrder {
			o.RestaurantStatus = RestaurantPreparingOrder
		} else {
			return fmt.Errorf("restaurant satus: %s", o.RestaurantStatus)
		}
	}
	if status == RestaurantPreparedOrder {
		if o.RestaurantStatus != RestaurantPreparedOrder {
			o.RestaurantStatus = RestaurantPreparedOrder
		}
	}

	return nil
}

func (o *Order) SetCourierStatus(status CourierStatus) error {
	if o.CourierStatus == "" {
		o.CourierStatus = status
	}

	if status == CourierAssigned {
		if o.CourierStatus != CourierAssigned && o.CourierStatus != "" {
			return fmt.Errorf("courier satus: %s", o.CourierStatus)
		} else {
			o.CourierStatus = CourierAssigned
		}
	}

	if status == CourierTookOrder {
		if o.CourierStatus != CourierTookOrder && o.CourierStatus != "" {
			return fmt.Errorf("courier satus: %s", o.CourierStatus)
		} else {
			o.CourierStatus = CourierTookOrder
		}
	}
	if status == CourierDelivering {
		if o.CourierStatus != CourierDelivered {
			o.CourierStatus = CourierDelivering
		} else {
			return fmt.Errorf("courier satus: %s", o.CourierStatus)
		}
	}
	if status == CourierDelivered {
		if o.CourierStatus != CourierDelivered {
			o.CourierStatus = CourierDelivered
		}
	}

	return nil
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
		PaymentID:        uuid.Nil,
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

func (v TransactionStatus) String() string { return string(v) }
func (v RestaurantStatus) String() string  { return string(v) }
func (v OrderStatus) String() string       { return string(v) }
func (v CourierStatus) String() string     { return string(v) }

const (
	Created   OrderStatus = "order.created"
	Canceled  OrderStatus = "order.canceled"
	Completed OrderStatus = "order.completed"
)

const (
	CourierAssigned   CourierStatus = "courier.order.assigned"
	CourierTookOrder  CourierStatus = "courier.order.taken"
	CourierDelivering CourierStatus = "courier.delivering"
	CourierDelivered  CourierStatus = "courier.delivered"
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
