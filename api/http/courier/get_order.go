package courier

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetOrder() gin.HandlerFunc {
	type Request struct {
		OrderID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request
		if err := c.BindUri(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order, err := h.courierRepository.GetOrder(c, uuid.MustParse(r.OrderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}
