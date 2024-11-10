package restaurant

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetRestaurant() gin.HandlerFunc {
	type Request struct {
		RestaurantID string `uri:"id" binding:"required,uuid"`
	}

	return func(c *gin.Context) {
		var r Request

		if err := c.BindUri(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		restaurant, err := h.restaurantRepository.GetRestaurant(c, uuid.MustParse(r.RestaurantID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, restaurant)
	}
}
