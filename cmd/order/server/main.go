package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-feast/resty-backend/config"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/internal/closer"
	"github.com/go-feast/resty-backend/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	//db = db.WithContext(ctx)

	service.RunService(ctx, serviceName, version, c, RegisterOrderServiceRoutes(db))
}

func RegisterOrderServiceRoutes(db *gorm.DB) func(r chi.Router) {
	return func(r chi.Router) {

	}
}
