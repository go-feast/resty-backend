package restaurant

import (
	"github.com/go-chi/chi/v5"
	apihttp "github.com/go-feast/resty-backend/api/http"
	"github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetMenu() func(w http.ResponseWriter, r *http.Request) {
	return getMenu(h.repository)
}

func getMenu(restaurantRepository restaurant.Repository) func(w http.ResponseWriter, r *http.Request) {
	// request is thd `id` in the path

	type Response struct {
		ID    uuid.UUID         `json:"id"`
		Meals []restaurant.Meal `json:"meals"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id  = uuid.MustParse(chi.URLParam(r, "id"))
			ctx = r.Context()
		)

		menu, err := restaurantRepository.GetMenu(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{
			ID:    id,
			Meals: menu,
		}

		apihttp.Encode(w, response)
	}
}
