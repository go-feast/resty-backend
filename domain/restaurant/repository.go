package restaurant

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateRestaurant(ctx context.Context, r Restaurant) error
	AppendMeals(ctx context.Context, rid uuid.UUID, ms ...Meal) error
	// SetMealForRestaurant
	// GetRestaurant
	GetRestaurant(ctx context.Context, rid uuid.UUID) (Restaurant, error)
	GetMenu(ctx context.Context, rid uuid.UUID) ([]Meal, error)
}
