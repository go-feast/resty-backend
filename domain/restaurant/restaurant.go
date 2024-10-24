package restaurant

import (
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
)

type Restaurant struct {
	ID          uuid.UUID
	Name        string
	Description string
	Location    geo.Location
	Meals       []Meal
}

func (r *Restaurant) AppendMeal(ms ...Meal) {
	r.Meals = append(r.Meals, ms...)
}

func NewRestaurant(
	name string,
	description string,
	location geo.Location,
) Restaurant {
	return Restaurant{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Location:    location,
		Meals:       []Meal{},
	}
}

type Meal struct {
	ID   uuid.UUID
	Name string
	// ingredients etc.
}

func NewMeal(name string) Meal {
	return Meal{
		ID:   uuid.New(),
		Name: name,
	}
}
