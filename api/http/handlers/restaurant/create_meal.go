package restaurant

import (
	"github.com/go-chi/chi/v5"
	apihttp "github.com/go-feast/resty-backend/api/http"
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) CreateMeals() func(w http.ResponseWriter, r *http.Request) {
	return createMeals(h.repository)
}

func createMeals(restaurantRepository restaurant.Repository) func(http.ResponseWriter, *http.Request) {
	type Request struct {
		RestaurantID uuid.UUID `json:"omitempty"` // mapped to path
		Name         []string  `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		stringID := chi.URLParam(r, "id")
		id := uuid.MustParse(stringID)

		req, err := apihttp.Decode[Request](r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.RestaurantID = id

		meals := make([]restaurant.Meal, len(req.Name))
		for i, name := range req.Name {
			meals[i] = restaurant.NewMeal(name)
		}

		err = restaurantRepository.AppendMeals(r.Context(), req.RestaurantID, meals...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
