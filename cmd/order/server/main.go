package main

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	httporder "github.com/go-feast/resty-backend/api/http/handlers/order"
	"github.com/go-feast/resty-backend/config"
	"github.com/go-feast/resty-backend/domain/order"
	outboxorder "github.com/go-feast/resty-backend/infrastructure/outbox/order"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/order/gorm"
	"github.com/go-feast/resty-backend/internal/closer"
	"github.com/go-feast/resty-backend/internal/logging"
	"github.com/go-feast/resty-backend/internal/pubsub"
	"github.com/go-feast/resty-backend/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os/signal"
	"syscall"
)

const (
	version     = "v1.0"
	serviceName = "order"
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

	order.InitializeOrderScheme(db)

	publisher, err := pubsub.NewSQLPublisher(db, logging.NewWatermillLogger())
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to database")
	}

	service.RunService(ctx, serviceName, version, c, RegisterOrderServiceRoutes(
		otel.GetTracerProvider().Tracer("order-service"),
		db,
		publisher,
	))
}

func RegisterOrderServiceRoutes(tracer trace.Tracer, db *gorm.DB, sqlPub message.Publisher) func(r chi.Router) {
	repository := gormorder.NewOrderRepository(db)
	outbox := outboxorder.NewOutbox(sqlPub, repository, json.Marshal)
	handler := httporder.NewHandler(
		tracer, repository, outbox)

	return func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

		r.Route("/api/v1", func(r chi.Router) {
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", handler.TakeOrder())
				r.Get("/{id}", handler.GetOrder())
			})
		})
	}
}
