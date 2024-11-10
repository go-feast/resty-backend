package order

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) TakeOrder() gin.HandlerFunc {
	type Request struct {
		CustomerID   uuid.UUID    `json:"customer_id"`
		RestaurantID uuid.UUID    `json:"restaurant_id"`
		Meals        uuid.UUIDs   `json:"meals"`
		Destination  geo.Location `json:"destination"`
	}

	return func(c *gin.Context) {
		r := &Request{}

		if err := c.MustBindWith(r, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		o := order.NewOrder(r.CustomerID, r.RestaurantID, r.Meals, r.Destination)

		err := h.orderRepository.Create(c, o)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := message.NewMessage(message.Event{"order_id": o.ID}, h.Marshaler)
		err = h.publisher.Publish(order.Created.String(), msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":        o.ID,
			"timestamp": o.CreatedAt,
		})
	}
}
