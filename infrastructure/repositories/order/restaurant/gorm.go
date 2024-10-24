package restaurant

import (
	"context"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitializeRestaurantModel(db *gorm.DB) {
	err := db.AutoMigrate(&order.Restaurant{})
	if err != nil {
		panic(err)
	}
}

type GormRestaurantRepository struct {
	db *gorm.DB
}

func NewGormRestaurantRepository(db *gorm.DB) *GormRestaurantRepository {
	return &GormRestaurantRepository{db: db}
}

func (g *GormRestaurantRepository) Create(ctx context.Context, restaurant *order.Restaurant) error {
	return g.db.WithContext(ctx).Create(restaurant).Error
}

func (g *GormRestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Restaurant, error) {
	r := &order.Restaurant{}
	tx := g.db.WithContext(ctx).Where("id = ?", id).First(r)
	if tx != nil {
		return nil, tx.Error
	}

	return r, nil
}
