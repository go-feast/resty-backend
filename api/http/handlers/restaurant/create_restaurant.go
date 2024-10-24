package restaurant

import (
	apihttp "github.com/go-feast/resty-backend/api/http"
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) CreateRestaurant() func(http.ResponseWriter, *http.Request) {
	return createRestaurant(h.repository)
}

func createRestaurant(
	restaurantRepository restaurant.Repository,
) func(http.ResponseWriter, *http.Request) {
	type Request struct {
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Location    geo.JSONLocation `json:"location"`
		Meals       []struct {
			Name string `json:"name"`
		}
	}

	type Response struct {
		ID uuid.UUID `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		req, err := apihttp.Decode[Request](r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := restaurant.NewRestaurant(
			req.Name,
			req.Description,
			geo.MustLocation(req.Location.Latitude, req.Location.Longitude),
		)
		meals := make([]restaurant.Meal, len(req.Meals))
		for i, meal := range req.Meals {
			meals[i] = restaurant.NewMeal(meal.Name)
		}

		res.AppendMeal(meals...)

		err = restaurantRepository.CreateRestaurant(ctx, res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{res.ID}

		apihttp.Encode(w, response)
	}
}
