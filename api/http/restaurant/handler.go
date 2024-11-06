package restaurant

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
)

type Handler struct {
	restaurantRepository restaurant.RestaurantRepository
	publisher            message.Publisher
	orderRepository      restaurant.OrderRepository
	Marshaler            func(interface{}) ([]byte, error)
	Unmarshaler          func([]byte, interface{}) error
}

func NewHandler(restaurantRepository restaurant.RestaurantRepository, publisher message.Publisher, orderRepository restaurant.OrderRepository) *Handler {
	return &Handler{restaurantRepository: restaurantRepository, publisher: publisher, orderRepository: orderRepository}
}
