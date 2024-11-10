package restaurant

import (
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
	"github.com/go-feast/resty-backend/internal/domain/shared/geo"
	"net/http"
)

func (h *Handler) CreateRestaurant() gin.HandlerFunc {
	type Request struct {
		Name     string       `json:"name"`
		Location geo.Location `json:"location"`
		Meals    []string
	}

	return func(c *gin.Context) {
		r := &Request{}

		if err := c.Bind(r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rest := restaurant.NewRestaurant(
			r.Name,
			r.Location,
			r.Meals,
		)

		err := h.restaurantRepository.CreateRestaurant(c, rest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, rest)
	}
}
