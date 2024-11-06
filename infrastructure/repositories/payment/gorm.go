package payment

import (
	"context"
	"github.com/go-feast/resty-backend/internal/domain/payment"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitializePaymentOrDie(db *gorm.DB) {
	err := db.AutoMigrate(&payment.Payment{})
	if err != nil {
		panic(err)
	}
}

type GormPaymentRepository struct {
	db *gorm.DB
}

func (g *GormPaymentRepository) Create(ctx context.Context, payment *payment.Payment) error {
	return g.db.WithContext(ctx).Create(payment).Error
}

func withTx(db *gorm.DB) *GormPaymentRepository {
	return &GormPaymentRepository{db: db}
}

func (g *GormPaymentRepository) Get(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	var p payment.Payment
	if err := g.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (g *GormPaymentRepository) Transact(ctx context.Context, id uuid.UUID, action func(*payment.Payment) error) (*payment.Payment, error) {
	var pay *payment.Payment
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		p, err := withTx(tx).Get(ctx, id)
		if err != nil {
			return err
		}

		err = action(p)
		if err != nil {
			return err
		}

		pay = p

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pay, nil
}

func NewGormPaymentRepository(db *gorm.DB) *GormPaymentRepository {
	return &GormPaymentRepository{db: db}
}
