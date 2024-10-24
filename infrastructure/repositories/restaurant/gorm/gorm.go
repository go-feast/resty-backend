package gorm

import (
	"context"
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitializeRestaurantScheme(db *gorm.DB) {
	err := db.AutoMigrate(&Meal{}, &Restaurant{})
	if err != nil {
		panic(err)
	}
}

type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (g *RestaurantRepository) GetRestaurant(ctx context.Context, rid uuid.UUID) (restaurant.Restaurant, error) {
	r := restaurant.Restaurant{}
	tx := g.db.WithContext(ctx).Where("id = ?", rid).First(&r)
	if tx.Error != nil {
		return restaurant.Restaurant{}, tx.Error
	}

	return r, nil
}

func (g *RestaurantRepository) GetMenu(ctx context.Context, rid uuid.UUID) ([]restaurant.Meal, error) {
	var r Restaurant
	tx := g.db.WithContext(ctx).Preload("Meals").Where("id = ?", rid).Find(&r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return toDomainMeals(r.Meals), nil
}

func (g *RestaurantRepository) AppendMeals(ctx context.Context, rid uuid.UUID, ms ...restaurant.Meal) error {
	var r = Restaurant{}
	tx := g.db.WithContext(ctx).Where("id = ?", rid).First(&r)
	if tx.Error != nil {
		return tx.Error
	}

	domainRestaurant := toDomainRestaurant(r)
	domainRestaurant.AppendMeal(ms...)

	r = mapRestaurant(domainRestaurant)

	return g.db.WithContext(ctx).Where("id = ?", rid).Save(r).Error
}

func (g *RestaurantRepository) CreateRestaurant(ctx context.Context, r restaurant.Restaurant) error {
	mapped := mapRestaurant(r)
	return g.db.WithContext(ctx).Create(&mapped).Error
}
