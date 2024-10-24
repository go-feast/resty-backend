package gorm

import (
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
)

type Restaurant struct {
	ID          uuid.UUID `order:"type:uuid;primaryKey"`
	Name        string    `order:"type:text"`
	Description string    `order:"type:text"`
	Latitude    float64   `order:"type:float"`
	Longitude   float64   `order:"type:float"`
	Meals       []Meal    `order:"foreignKey:RestaurantID;references:ID"`
}

type Meal struct {
	ID           uuid.UUID `order:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID `order:"type:uuid"`
	Name         string    `order:"type:text"`
	// ingredients etc.
}

func mapRestaurant(restaurant restaurant.Restaurant) Restaurant {
	return Restaurant{
		ID:          restaurant.ID,
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Latitude:    restaurant.Location.Latitude(),
		Longitude:   restaurant.Location.Longitude(),
		Meals:       mapMeals(restaurant.ID, restaurant.Meals),
	}
}

func mapMeals(rid uuid.UUID, meals []restaurant.Meal) []Meal {
	res := make([]Meal, len(meals))
	for i, meal := range meals {
		res[i] = Meal{
			ID:           meal.ID,
			RestaurantID: rid,
			Name:         meal.Name,
		}
	}
	return res
}

func toDomainRestaurant(r Restaurant) restaurant.Restaurant {
	destination, _ := geo.NewDestination(r.Latitude, r.Longitude)
	return restaurant.Restaurant{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Location:    destination,
		Meals:       nil,
	}
}

func toDomainMeals(meals []Meal) []restaurant.Meal {
	res := make([]restaurant.Meal, len(meals))
	for i, meal := range meals {
		res[i] = restaurant.Meal{
			ID:   meal.ID,
			Name: meal.Name,
		}
	}
	return res
}
