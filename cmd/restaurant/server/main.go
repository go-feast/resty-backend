package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	httprestaurant "github.com/go-feast/resty-backend/api/http/handlers/restaurant"
	"github.com/go-feast/resty-backend/config"
	gormrestaurant "github.com/go-feast/resty-backend/infrastructure/repositories/restaurant/gorm"
	"github.com/go-feast/resty-backend/internal/closer"
	"github.com/go-feast/resty-backend/internal/service"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os/signal"
	"syscall"
)

const (
	version     = "v1.0.0"
	serviceName = "restaurant"
)

func main() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		forClose    = closer.NewCloser()
		c           = config.ServiceConfig{}
	)

	defer forClose.Close()
	defer cancel()

	log.Logger = log.Logger.With().
		Str("app", serviceName+version).Caller().Timestamp().
		Logger()

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Panic().Err(err).Msg("failed to parse config")
	}

	log.Info().Any("config", c).Send()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           "pgx/v5",
		DSN:                  c.DB.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to database")
	}

	gormrestaurant.InitializeRestaurantScheme(db)

	service.RunService(ctx, serviceName, version, c, RegisterOrderServiceRoutes(
		db,
	))
}

func RegisterOrderServiceRoutes(db *gorm.DB) func(r chi.Router) {
	repository := gormrestaurant.NewRestaurantRepository(db)
	handler := httprestaurant.NewHandler(repository)

	return func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

		r.Route("/api/v1", func(r chi.Router) {
			r.Route("/restaurants", func(r chi.Router) {
				r.Post("/", handler.CreateRestaurant())
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handler.GetRestaurant())
					r.Route("/meals", func(r chi.Router) {
						r.Get("/", handler.GetMenu())
						r.Post("/", handler.CreateMeals())
					})
				})
			})
		})
	}
}
