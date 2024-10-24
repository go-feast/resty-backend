package order

import (
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
	"time"
)

type OrderFactory interface {
	Build() *Order

	WithID(id uuid.UUID) OrderFactory
	WithRestaurantID(id uuid.UUID) OrderFactory
	WithCustomerID(id uuid.UUID) OrderFactory
	WithCourierID(id uuid.UUID) OrderFactory
	WithMeals(ms uuid.UUIDs) OrderFactory
	WithState(state string) OrderFactory
	WithTransactionID(id uuid.UUID) OrderFactory
	WithDestination(loc geo.Location) OrderFactory
	WithCreatedAt(t time.Time) OrderFactory
}

type orderFactory struct {
	id            uuid.UUID
	restaurantID  uuid.UUID
	customerID    uuid.UUID
	courierID     uuid.UUID
	meals         uuid.UUIDs
	state         State
	transactionID uuid.UUID
	destination   geo.Location
	createdAt     time.Time
}

func (o *orderFactory) Build() *Order {
	return &Order{
		id:            o.id,
		restaurantID:  o.restaurantID,
		customerID:    o.customerID,
		courierID:     o.courierID,
		meals:         o.meals,
		state:         o.state,
		transactionID: o.transactionID,
		destination:   o.destination,
		createdAt:     o.createdAt,
	}
}

func (o *orderFactory) WithID(id uuid.UUID) OrderFactory {
	o.id = id
	return o
}

func (o *orderFactory) WithRestaurantID(id uuid.UUID) OrderFactory {
	o.restaurantID = id
	return o
}

func (o *orderFactory) WithCustomerID(id uuid.UUID) OrderFactory {
	o.customerID = id
	return o
}

func (o *orderFactory) WithCourierID(id uuid.UUID) OrderFactory {
	o.courierID = id
	return o
}

func (o *orderFactory) WithMeals(ms uuid.UUIDs) OrderFactory {
	o.meals = ms
	return o
}

func (o *orderFactory) WithState(state string) OrderFactory {
	s, err := GetState(state)
	if err != nil {
		panic(err)
	}
	o.state = *s
	return o
}

func (o *orderFactory) WithTransactionID(id uuid.UUID) OrderFactory {
	o.transactionID = id
	return o
}

func (o *orderFactory) WithDestination(loc geo.Location) OrderFactory {
	o.destination = loc
	return o
}

func (o *orderFactory) WithCreatedAt(t time.Time) OrderFactory {
	o.createdAt = t
	return o
}

func NewOrderFactory() OrderFactory {
	return &orderFactory{}
}
