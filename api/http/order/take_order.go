package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func takeOrder() gin.HandlerFunc {
	type Request struct {
		CustomerID   uuid.UUID
		RestaurantID uuid.UUID
		Meals        uuid.UUIDs
	}

	return func(c *gin.Context) {

	}
}
