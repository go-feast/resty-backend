package restaurnt

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

func InitializeRestaurantOrDie(db *gorm.DB) {
	err := db.AutoMigrate(&restaurant.Restaurant{}, &restaurant.Meal{}, &restaurant.Order{})
	if err != nil {
		log.Fatal(err)
	}
}

type GormRestaurantRepository struct {
	db *gorm.DB
}

func NewGormRestaurantRepository(db *gorm.DB) *GormRestaurantRepository {
	return &GormRestaurantRepository{db: db}
}

func (g *GormRestaurantRepository) GetRestaurant(ctx context.Context, id uuid.UUID) (*restaurant.Restaurant, error) {
	r := &restaurant.Restaurant{}
	result := g.db.WithContext(ctx).Preload("Meals").Where("id = ?", id).First(r)
	if result.Error != nil {
		return nil, result.Error
	}

	return r, nil
}

func withTx(db *gorm.DB) *GormRestaurantRepository {
	return &GormRestaurantRepository{db: db}
}

func (g *GormRestaurantRepository) Transact(ctx context.Context, id uuid.UUID, action func(restaurant *restaurant.Restaurant) error) (*restaurant.Restaurant, error) {
	var res *restaurant.Restaurant
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r, err := withTx(tx).GetRestaurant(ctx, id)
		if err != nil {
			return err
		}

		err = action(r)
		if err != nil {
			return err
		}

		tx.Save(r)

		res = r

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *GormRestaurantRepository) CreateRestaurant(ctx context.Context, restaurant *restaurant.Restaurant) error {
	return g.db.WithContext(ctx).Create(restaurant).Error
}
