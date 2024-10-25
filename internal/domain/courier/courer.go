package courier

import (
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"time"
)

type Courier struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	AssignedOrder Order
}

type Order struct {
	ID                  uuid.UUID
	RestaurantLocation  geo.Location
	Status              Status
	DestinationLocation geo.Location
	CreatedAt           time.Time
}

type Status string

const (
	Assigned   Status = "courier.order.assigned"
	TookOrder  Status = "courier.order.taken"
	Delivering Status = "courier.delivering"
	Delivered  Status = "courier.delivered"
)
