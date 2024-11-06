package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type (
	PaymentRepository interface {
		Create(ctx context.Context, payment *Payment) error
		Get(ctx context.Context, id uuid.UUID) (*Payment, error)
		Transact(ctx context.Context, id uuid.UUID, action func(*Payment) error) (*Payment, error)
	}

	Payment struct {
		ID      uuid.UUID
		OrderID uuid.UUID
		// Amount

		OrderStatus   OrderStatus
		PaymentStatus PaymentStatus
		CreatedAt     time.Time
		PaidAt        time.Time
	}

	PaymentStatus string
	OrderStatus   string
)

func NewPayment(orderID uuid.UUID) *Payment {
	return &Payment{
		ID:            uuid.New(),
		OrderID:       orderID,
		OrderStatus:   OrderCreated,
		PaymentStatus: Waiting,
		CreatedAt:     time.Now(),
		PaidAt:        time.Time{},
	}
}

func (p PaymentStatus) String() string { return string(p) }
func (p OrderStatus) String() string   { return string(p) }

func (p *Payment) IsOrderCanceledOrCompleted() bool {
	return p.OrderStatus == OrderCanceled || p.OrderStatus == OrderCompleted
}

func (p *Payment) SetPaymentStatus(status PaymentStatus) error {
	if status == Waiting {
		if p.PaymentStatus != Waiting && p.IsOrderCanceledOrCompleted() {
			return fmt.Errorf("payment status: %s; order status: %s", p.PaymentStatus, p.OrderStatus)
		}
	}
	if status == Paid {
		if p.PaymentStatus != Canceled && !p.IsOrderCanceledOrCompleted() {
			p.PaymentStatus = Paid
		} else {
			return fmt.Errorf("payment status: %s; order status: %s", p.PaymentStatus, p.OrderStatus)
		}
	}
	if status == Canceled {
		if p.PaymentStatus != Paid && !p.IsOrderCanceledOrCompleted() {
			p.PaymentStatus = Canceled
		} else {
			return fmt.Errorf("payment status: %s; order status: %s", p.PaymentStatus, p.OrderStatus)
		}
	}

	return nil
}

func (p *Payment) SetOrderStatus(status OrderStatus) error {
	if status == OrderCreated {
		if p.OrderStatus != OrderCreated {
			return fmt.Errorf("order satus: %s", p.OrderStatus)
		}
	}
	if status == OrderCanceled {
		if p.OrderStatus != OrderCompleted {
			p.OrderStatus = OrderCanceled
		} else {
			return fmt.Errorf("order satus: %s", p.OrderStatus)
		}
	}
	if status == OrderCompleted {
		p.OrderStatus = OrderCompleted
	}

	return nil

}

const (
	Waiting  PaymentStatus = "payment.waiting"
	Paid     PaymentStatus = "payment.paid"
	Canceled PaymentStatus = "payment.canceled"
)

const (
	OrderCreated   = "order.created"
	OrderCanceled  = "order.canceled"
	OrderCompleted = "order.completed"
)
