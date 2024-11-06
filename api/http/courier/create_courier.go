package courier

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/courier"
	"net/http"
)

func (h *Handler) CreateCourier() gin.HandlerFunc {
	type Request struct {
		Name string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindJSON(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cour := courier.NewCourier(r.Name)

		err := h.courierRepository.Create(c.Request.Context(), cour)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, cour)
	}
}
