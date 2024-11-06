package restaurant

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) SetPreparingOrder() gin.HandlerFunc {
	type Request struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindUri(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		o, err := h.orderRepository.Transact(c, uuid.MustParse(r.OrderID), func(order *restaurant.Order) error {
			return order.SetRestaurantStatus(restaurant.PreparingOrder)
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, o)
	}
}
