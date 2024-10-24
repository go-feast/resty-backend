package restaurant

import (
	domain "github.com/go-feast/resty-backend/domain/restaurant"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	repository domain.Repository

	// metrics
	menuRequested prometheus.CounterVec
}

func NewHandler(repository domain.Repository) *Handler {
	return &Handler{
		repository: repository,
		//menuRequested: promauto.NewCounterVec()
	}
}
