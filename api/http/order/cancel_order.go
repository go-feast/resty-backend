package order

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) CancelOrder() gin.HandlerFunc {
	type Request struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request
		if err := c.BindUri(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		o, err := h.orderRepository.Transact(c, uuid.MustParse(r.OrderID), func(o *order.Order) error {
			o.OrderStatus = order.Canceled
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := message.NewMessage(message.Event{"order_id": o.ID}, h.Marshaler)
		err = h.publisher.Publish(order.Canceled.String(), msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"action": "canceled", "order": o})
	}
}
