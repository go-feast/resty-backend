package order

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"gorm.io/gorm"
	"log"
)

func InitializerOrderOrDie(db *gorm.DB) {
	err := db.AutoMigrate(&order.Meal{}, &order.Order{})
	if err != nil {
		log.Fatal(err)
	}
}

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (g *GormOrderRepository) Create(ctx context.Context, order *order.Order) error {
	return g.db.WithContext(ctx).Create(order).Error
}
