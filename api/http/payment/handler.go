package payment

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/payment"
)

type Handler struct {
	paymentRepository payment.PaymentRepository
	publisher         message.Publisher
	Unmarshaler       func([]byte, interface{}) error
	Marshaler         func(interface{}) ([]byte, error)
}

func NewHandler(paymentRepository payment.PaymentRepository, publisher message.Publisher, unmarshaler func([]byte, interface{}) error, marshaler func(interface{}) ([]byte, error)) *Handler {
	return &Handler{paymentRepository: paymentRepository, publisher: publisher, Unmarshaler: unmarshaler, Marshaler: marshaler}
}
