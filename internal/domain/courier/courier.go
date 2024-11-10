package courier

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"time"
)

type CourierRepository interface {
	Create(ctx context.Context, courier *Courier) error
	CreateOrder(ctx context.Context, order *Order) error
	Get(ctx context.Context, id uuid.UUID) (*Courier, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	AssignOrder(ctx context.Context, cid, oid uuid.UUID) error
	Transact(ctx context.Context, oid uuid.UUID, f func(*Order) error) (*Order, error)
}

type Courier struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Name           string
	AssignedOrders []Order
}

func NewCourier(name string) *Courier {
	return &Courier{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		Name:           name,
		AssignedOrders: nil,
	}
}

type Order struct {
	ID                           uuid.UUID
	CourierID                    *uuid.UUID
	RestaurantLocationLatitude   float64
	RestaurantLocationLongitude  float64
	Status                       Status
	DestinationLocationLatitude  float64
	DestinationLocationLongitude float64
	CreatedAt                    time.Time
}

func NewOrder(
	ID uuid.UUID,
	restaurantLocation geo.Location,
	destinationLocation geo.Location,
	createdAt time.Time,
) *Order {
	return &Order{
		ID:                           ID,
		CourierID:                    nil,
		RestaurantLocationLatitude:   restaurantLocation.Latitude,
		RestaurantLocationLongitude:  restaurantLocation.Longitude,
		Status:                       "",
		DestinationLocationLatitude:  destinationLocation.Latitude,
		DestinationLocationLongitude: destinationLocation.Longitude,
		CreatedAt:                    createdAt,
	}
}

type Status string

func (s Status) String() string { return string(s) }

const (
	Assigned   Status = "courier.order.assigned"
	TookOrder  Status = "courier.order.taken"
	Delivering Status = "courier.delivering"
	Delivered  Status = "courier.delivered"
)
