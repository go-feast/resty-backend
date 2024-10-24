package restaurant

import (
	"github.com/go-chi/chi/v5"
	apihttp "github.com/go-feast/resty-backend/api/http"
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetRestaurant() func(w http.ResponseWriter, r *http.Request) {
	return getRestaurant(h.repository)
}

func getRestaurant(restaurantRepository restaurant.Repository) func(http.ResponseWriter, *http.Request) {
	type Response struct {
		ID          uuid.UUID         `json:"id"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Location    geo.JSONLocation  `json:"location"`
		Menu        []restaurant.Meal `json:"menu"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.MustParse(chi.URLParam(r, "id"))

		res, err := restaurantRepository.GetRestaurant(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{
			ID:          res.ID,
			Name:        res.Name,
			Description: res.Description,
			Location:    res.Location.ToJSON(),
			Menu:        res.Meals,
		}

		apihttp.Encode(w, response)
	}
}
