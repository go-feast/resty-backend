package order

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func withTx(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (g *GormRepository) CreateOrder(ctx context.Context, order *restaurant.Order) error {
	return g.db.WithContext(ctx).Create(order).Error
}

func (g *GormRepository) GetOrder(ctx context.Context, id uuid.UUID) (*restaurant.Order, error) {
	var order restaurant.Order
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (g *GormRepository) Transact(ctx context.Context, id uuid.UUID, action func(order *restaurant.Order) error) (*restaurant.Order, error) {
	var o *restaurant.Order
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := withTx(tx).GetOrder(ctx, id)
		if err != nil {
			return err
		}

		err = action(order)
		if err != nil {
			return err
		}

		tx.Model(&restaurant.Order{}).Where("id = ?", id).Updates(order)

		o = order
		return nil
	})
	if err != nil {
		return nil, err
	}

	return o, nil
}
