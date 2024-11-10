package courier

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/courier"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func InitializeCourierOrDie(db *gorm.DB) {
	err := db.AutoMigrate(&courier.Courier{}, &courier.Order{})
	if err != nil {
		panic(err)
	}
}

type GormCourierRepository struct {
	db *gorm.DB
}

func (g *GormCourierRepository) GetOrder(ctx context.Context, id uuid.UUID) (*courier.Order, error) {
	var o courier.Order
	err := g.db.WithContext(ctx).Model(&courier.Order{}).Where("id = ?", id).First(&o)
	if err != nil {
		return nil, err.Error
	}

	return &o, nil
}

func (g *GormCourierRepository) CreateOrder(ctx context.Context, order *courier.Order) error {
	return g.db.WithContext(ctx).Create(order).Error
}

func NewGormCourierRepository(db *gorm.DB) *GormCourierRepository {
	return &GormCourierRepository{db: db}
}

func (g *GormCourierRepository) Create(ctx context.Context, courier *courier.Courier) error {
	return g.db.WithContext(ctx).Create(courier).Error
}

func (g *GormCourierRepository) Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error) {
	var c courier.Courier
	if err := g.db.WithContext(ctx).Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (g *GormCourierRepository) AssignOrder(ctx context.Context, cid, oid uuid.UUID) error {
	var c courier.Courier
	err := g.db.WithContext(ctx).Preload("AssignedOrders").Where("id = ?", cid).First(&c)
	if err.Error != nil {
		return errors.Wrap(err.Error, "failed to find courier")
	}

	var o courier.Order
	err = g.db.WithContext(ctx).Model(&courier.Order{}).Where("id = ?", oid).First(&o)
	if err.Error != nil {
		return errors.Wrap(err.Error, "failed to find order")
	}

	o.CourierID = &cid

	c.AssignedOrders = append(c.AssignedOrders, o)

	return g.db.WithContext(ctx).Save(&c).Error
}

func (g *GormCourierRepository) Transact(ctx context.Context, oid uuid.UUID, f func(*courier.Order) error) (*courier.Order, error) {
	var res *courier.Order
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var o courier.Order
		t := g.db.Model(&courier.Order{}).Where("id = ?", oid).First(&o)
		if t.Error != nil {
			return t.Error
		}

		err := f(&o)
		if err != nil {
			return err
		}

		res = &o

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
