package courier

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/courier"
	"github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) DeliveringOrder() gin.HandlerFunc {
	type Request struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindUri(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		o, err := h.courierRepository.Transact(c, uuid.MustParse(r.OrderID), func(order *courier.Order) error {
			order.Status = courier.Delivering
			return nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := message.NewMessage(message.Event{"order_id": r.OrderID}, json.Marshal)
		err = h.publisher.Publish(courier.Delivering.String(), msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, o)
	}
}
