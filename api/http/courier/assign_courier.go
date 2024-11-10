package courier

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/courier"
	"github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) AssignCourier() gin.HandlerFunc {
	type Request struct {
		CourierID string `uri:"cid" binding:"required,uuid"`
		OrderID   string `uri:"oid" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindUri(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.courierRepository.AssignOrder(c, uuid.MustParse(r.CourierID), uuid.MustParse(r.OrderID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := message.NewMessage(message.Event{"order_id": r.OrderID, "courier_id": r.CourierID}, json.Marshal)
		err = h.publisher.Publish(courier.Assigned.String(), msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
