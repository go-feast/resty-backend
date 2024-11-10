package courier

import "github.com/go-feast/resty-backend/internal/domain/courier"

type Handler struct {
	courierRepository courier.CourierRepository
	Marshaler         func(interface{}) ([]byte, error)
	Unmarshaler       func([]byte, interface{}) error
}

func NewHandler(courierRepository courier.CourierRepository, marshaler func(interface{}) ([]byte, error), unmarshaler func([]byte, interface{}) error) *Handler {
	return &Handler{courierRepository: courierRepository, Marshaler: marshaler, Unmarshaler: unmarshaler}
}
