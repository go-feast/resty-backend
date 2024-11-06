package order

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
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

func withTx(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (g *GormOrderRepository) Transact(ctx context.Context, id uuid.UUID, action func(o *order.Order) error) (*order.Order, error) {
	var (
		o   *order.Order
		err error
	)

	tx := g.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	o, err = withTx(tx).GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	err = action(o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (g *GormOrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	var o order.Order
	tx := g.db.Preload("Meals").Where("id = ?", id).First(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &o, nil
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (g *GormOrderRepository) Create(ctx context.Context, order *order.Order) error {
	return g.db.WithContext(ctx).Create(order).Error
}
