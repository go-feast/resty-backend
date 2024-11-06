package courier

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type CourierRepository interface {
	Create(ctx context.Context, courier *Courier) error
	Get(ctx context.Context, id uuid.UUID) (*Courier, error)
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
	CourierID                    uuid.UUID
	RestaurantLocationLatitude   float64
	RestaurantLocationLongitude  float64
	Status                       Status
	DestinationLocationLatitude  float64
	DestinationLocationLongitude float64
	CreatedAt                    time.Time
}

type Status string

func (s Status) String() string { return string(s) }

const (
	Assigned   Status = "courier.order.assigned"
	TookOrder  Status = "courier.order.taken"
	Delivering Status = "courier.delivering"
	Delivered  Status = "courier.delivered"
)
