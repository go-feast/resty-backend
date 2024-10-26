package order

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) TakeOrder() gin.HandlerFunc {
	return takeOrder(h.orderRepository)
}

func takeOrder(orderRepository order.OrderRepository) gin.HandlerFunc {
	type Request struct {
		CustomerID   uuid.UUID
		RestaurantID uuid.UUID
		Meals        uuid.UUIDs
		Destination  geo.Location
	}

	return func(c *gin.Context) {
		r := &Request{}

		if err := c.MustBindWith(r, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		o := order.NewOrder(r.CustomerID, r.RestaurantID, r.Meals, r.Destination)

		err := orderRepository.Create(c, o)
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
