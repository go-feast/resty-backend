package order

import (
	"context"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) Delete(ctx context.Context, o *order.Order) error {
	result := r.db.WithContext(ctx).Delete(toDatabaseDTO(o))
	if result.Error != nil {
		return errors.Wrap(result.Error, "order repository: failed to delete order")
	}

	return nil
}

func NewOrderRepository(
	db *gorm.DB,
) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	result := r.db.WithContext(ctx).Create(toDatabaseDTO(o))
	if result.Error != nil {
		return errors.Wrap(result.Error, "order repository: failed to create order")
	}

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o := &DatabaseOrderDTO{}

	result := r.db.WithContext(ctx).Find(o, "id = ?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "order repository: order get: failed to find order")
	}

	return toOrder(o), nil
}

func (r *OrderRepository) Operate(ctx context.Context, id uuid.UUID, op order.Operation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		// select order to escape data race
		o, err := withTx(tx).Get(ctx, id)
		if err != nil {
			return errors.Wrap(err, "order operate: failed to get order")
		}

		err = op(o)
		if err != nil {
			return errors.Wrap(err, "order operate: failed to operate order")
		}

		result := tx.Save(toDatabaseDTO(o))
		if result.Error != nil {
			return errors.Wrap(result.Error, "order operate: failed to save order")
		}

		return nil
	})
}

func withTx(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{db: tx}
}
