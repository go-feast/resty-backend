package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/payment"
	"github.com/go-feast/resty-backend/internal/message"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) PayForOrder() gin.HandlerFunc {
	type Request struct {
		PaymentID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindUri(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p, err := h.paymentRepository.Transact(c, uuid.MustParse(r.PaymentID), func(p *payment.Payment) error {
			return p.SetPaymentStatus(payment.Paid)
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := message.NewMessage(message.Event{"order_id": p.OrderID}, h.Marshaler)
		err = h.publisher.Publish(payment.Paid.String(), msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}
