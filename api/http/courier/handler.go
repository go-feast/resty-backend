package courier

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/courier"
)

type Handler struct {
	courierRepository courier.CourierRepository
	publisher         message.Publisher
}

func NewHandler(courierRepository courier.CourierRepository, publisher message.Publisher) *Handler {
	return &Handler{courierRepository: courierRepository, publisher: publisher}
}
