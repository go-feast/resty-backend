package order

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetOrder() gin.HandlerFunc {
	return getOrder(h.orderRepository)
}

func getOrder(orderRepository order.OrderRepository) gin.HandlerFunc {
	type Request struct {
		ID string `uri:"id" binding:"required,uuid"`
	}
	return func(c *gin.Context) {
		r := &Request{}
		if err := c.BindUri(r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		o, err := orderRepository.GetOrder(c.Request.Context(), uuid.MustParse(r.ID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, o)
	}
}
