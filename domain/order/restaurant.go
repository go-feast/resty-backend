package order

import (
	"context"
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
)

type Restaurant struct {
	ID       uuid.UUID
	Location geo.Location
}

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *Restaurant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Restaurant, error)
}
