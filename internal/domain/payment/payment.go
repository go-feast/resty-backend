package payment

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	// Amount
	Status    Status
	CreatedAt time.Time
	PaidAt    time.Time
}

type Status string

const (
	Waiting  Status = "payment.waiting"
	Paid     Status = "payment.paid"
	Canceled Status = "payment.canceled"
)
