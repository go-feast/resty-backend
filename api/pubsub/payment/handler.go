package payment

import "github.com/go-feast/resty-backend/internal/domain/payment"

type Handler struct {
	paymentRepository payment.PaymentRepository
	Unmarshaler       func([]byte, interface{}) error
	Marshaler         func(interface{}) ([]byte, error)
}

func NewHandler(paymentRepository payment.PaymentRepository, unmarshaler func([]byte, interface{}) error, marshaler func(interface{}) ([]byte, error)) *Handler {
	return &Handler{paymentRepository: paymentRepository, Unmarshaler: unmarshaler, Marshaler: marshaler}
}
